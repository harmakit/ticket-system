package service

import (
	"context"
	"ticket-system/checkout/internal/client/booking"
	"ticket-system/checkout/internal/client/event"
	"ticket-system/checkout/internal/model"
	"ticket-system/lib/query-engine/postgres"
)

type OrderService interface {
	GetOrder(ctx context.Context, id model.UUID) (*model.Order, error)
	CreateOrder(ctx context.Context, o *model.Order) error
	UpdateOrder(ctx context.Context, order *model.Order) error
}

type ItemService interface {
	GetItem(ctx context.Context, id model.UUID) (*model.Item, error)
	ListItems(ctx context.Context, orderId model.UUID) ([]*model.Item, error)
	CreateItem(ctx context.Context, i *model.Item) error
}

type CartService interface {
	GetCart(ctx context.Context, id model.UUID) (*model.Cart, error)
	ListCarts(ctx context.Context, userId *model.UUID, ticketId *model.UUID) ([]*model.Cart, error)
	CreateCart(ctx context.Context, c *model.Cart) error
	UpdateCart(ctx context.Context, c *model.Cart) error
	DeleteCart(ctx context.Context, c *model.Cart) error
}

type BusinessLogic struct {
	transactionManager *postgres.TransactionManager

	orderService OrderService
	itemService  ItemService
	cartService  CartService

	bookingService booking.Service
	eventService   event.Service
}

func New(
	transactionManager *postgres.TransactionManager,

	orderService OrderService,
	itemService ItemService,
	cartService CartService,

	bookingService booking.Service,
	eventService event.Service,
) BusinessLogic {
	return BusinessLogic{
		transactionManager,
		orderService,
		itemService,
		cartService,
		bookingService,
		eventService,
	}
}

func (s BusinessLogic) GetOrder(ctx context.Context, id model.UUID) (*model.Order, []*model.Item, error) {
	order, err := s.orderService.GetOrder(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	items, err := s.itemService.ListItems(ctx, order.Id)
	if err != nil {
		return nil, nil, err
	}

	return order, items, nil
}

func (s BusinessLogic) CreateOrder(ctx context.Context, o *model.Order) error {
	return s.orderService.CreateOrder(ctx, o)
}

func (s BusinessLogic) GetItem(ctx context.Context, id model.UUID) (*model.Item, error) {
	return s.itemService.GetItem(ctx, id)
}

func (s BusinessLogic) CreateItem(ctx context.Context, i *model.Item) error {
	return s.itemService.CreateItem(ctx, i)
}

func (s BusinessLogic) GetCart(ctx context.Context, id model.UUID) (*model.Cart, error) {
	return s.cartService.GetCart(ctx, id)
}

func (s BusinessLogic) AddToCart(ctx context.Context, c *model.Cart) error {
	carts, err := s.cartService.ListCarts(ctx, &c.UserId, &c.TicketId)
	if err != nil {
		return err
	}

	if len(carts) == 0 {
		return s.cartService.CreateCart(ctx, c)
	}

	if len(carts) > 1 {
		return ErrMultipleCartsReceived
	}

	c.Id = carts[0].Id

	return s.cartService.UpdateCart(ctx, c)
}

func (s BusinessLogic) UpdateCart(ctx context.Context, c *model.Cart) error {
	return s.cartService.UpdateCart(ctx, c)
}

func (s BusinessLogic) PlaceOrder(ctx context.Context, userId model.UUID) error {
	carts, err := s.cartService.ListCarts(ctx, &userId, nil)
	if err != nil {
		return err
	}
	cartsLen := len(carts)

	if cartsLen == 0 {
		return ErrEmptyCart
	}

	tickets := make([]*model.Ticket, cartsLen)
	stocks := make([]*model.Stock, cartsLen)
	bookingsIds := make([]model.UUID, cartsLen)

	for i, cart := range carts {
		receivedTicket, err := s.eventService.GetTicket(ctx, cart.TicketId)
		if err != nil {
			return err
		}
		tickets[i] = receivedTicket

		receivedStock, err := s.bookingService.GetTicketStock(ctx, receivedTicket.EventId, receivedTicket.Id)
		if err != nil {
			return err
		}

		stocks[i] = receivedStock
	}

	order := &model.Order{
		UserId: userId,
		Status: model.StatusCreated,
	}

	err = s.transactionManager.RunReadCommitted(ctx, func(ctx context.Context) error {
		err = s.orderService.CreateOrder(ctx, order)
		if err != nil {
			return err
		}

		items := make([]*model.Item, cartsLen)
		for i := range cartsLen {
			bookingId, err := s.bookingService.CreateBooking(ctx, tickets[i], order, carts[i], userId)
			if err != nil {
				return err
			}
			bookingsIds[i] = bookingId

			items[i] = &model.Item{
				OrderId: order.Id,
				StockId: stocks[i].Id,
				Count:   carts[i].Count,
			}

			err = s.itemService.CreateItem(ctx, items[i])
			if err != nil {
				return err
			}

			err = s.cartService.DeleteCart(ctx, carts[i])
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return s.bookingService.ExpireBookings(ctx, bookingsIds)
	}

	return nil
}

func (s BusinessLogic) MarkOrderAsPaid(ctx context.Context, id model.UUID) error {
	order, err := s.orderService.GetOrder(ctx, id)
	if err != nil {
		return err
	}

	err = s.bookingService.DeleteOrderBookings(ctx, order.Id)
	if err != nil {
		return err
	}

	order.Status = model.StatusPaid
	return s.orderService.UpdateOrder(ctx, order)
}
