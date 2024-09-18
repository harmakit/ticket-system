package v1

import (
	"context"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	"ticket-system/checkout/internal/model"
	"ticket-system/checkout/internal/service"
	desc "ticket-system/checkout/pkg/v1"
)

type CheckoutServiceImplementation struct {
	desc.UnsafeCheckoutServiceServer
	bl service.BusinessLogic
}

func NewCheckoutServiceImplementation(bl service.BusinessLogic) desc.CheckoutServiceServer {
	return &CheckoutServiceImplementation{bl: bl}
}

func (impl CheckoutServiceImplementation) GetOrder(ctx context.Context, req *desc.GetOrderRequest) (*desc.GetOrderResponse, error) {
	var res desc.GetOrderResponse

	o, is, err := impl.bl.GetOrder(ctx, model.UUID(req.Id))
	if err != nil {
		return nil, err
	}

	res.Order = impl.bindModelToDescOrder(o, is)

	return &res, nil
}

func (impl CheckoutServiceImplementation) ListOrders(ctx context.Context, req *desc.ListOrdersRequest) (*desc.ListOrdersResponse, error) {
	var res desc.ListOrdersResponse

	orders, items, err := impl.bl.ListOrders(ctx, model.UUID(req.UserId), int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, err
	}

	resOrders := make([]*desc.Order, 0, len(orders))
	for i, order := range orders {
		resOrders = append(resOrders, impl.bindModelToDescOrder(order, items[i]))
	}

	res = desc.ListOrdersResponse{
		Orders: resOrders,
	}

	return &res, nil
}

func (impl CheckoutServiceImplementation) AddToCart(ctx context.Context, req *desc.AddToCartRequest) (*desc.AddToCartResponse, error) {
	var res desc.AddToCartResponse

	c := &model.Cart{
		UserId:   model.UUID(req.UserId),
		TicketId: model.UUID(req.TicketId),
		Count:    int(req.Count),
	}
	err := impl.bl.AddToCart(ctx, c)
	if err != nil {
		return nil, err
	}

	res = desc.AddToCartResponse{
		Id:       string(c.Id),
		UserId:   string(c.UserId),
		TicketId: string(c.TicketId),
		Count:    int32(c.Count),
	}

	return &res, nil
}

func (impl CheckoutServiceImplementation) UpdateCart(ctx context.Context, req *desc.UpdateCartRequest) (*emptypb.Empty, error) {
	var res emptypb.Empty

	c := &model.Cart{
		Id:       model.UUID(req.Id),
		UserId:   model.UUID(req.UserId),
		TicketId: model.UUID(req.TicketId),
		Count:    int(req.Count),
	}
	err := impl.bl.UpdateCart(ctx, c)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (impl CheckoutServiceImplementation) GetUserCart(ctx context.Context, req *desc.GetUserCartRequest) (*desc.GetUserCartResponse, error) {
	var res desc.GetUserCartResponse

	carts, err := impl.bl.GetUserCart(ctx, model.UUID(req.UserId))
	if err != nil {
		return nil, err
	}

	items := make([]*desc.Cart, 0, len(carts))
	for _, cart := range carts {
		items = append(items, impl.bindModelToDescCart(cart))
	}

	res = desc.GetUserCartResponse{
		Carts: items,
	}

	return &res, nil
}

func (impl CheckoutServiceImplementation) PlaceOrder(ctx context.Context, req *desc.PlaceOrderRequest) (*desc.PlaceOrderResponse, error) {
	var res desc.PlaceOrderResponse

	o, err := impl.bl.PlaceOrder(ctx, model.UUID(req.UserId))
	if err != nil {
		return nil, err
	}

	o, is, err := impl.bl.GetOrder(ctx, o.Id)
	if err != nil {
		return nil, err
	}

	res.Order = impl.bindModelToDescOrder(o, is)

	return &res, nil
}

func (impl CheckoutServiceImplementation) MarkOrderAsPaid(ctx context.Context, req *desc.MarkOrderAsPaidRequest) (*emptypb.Empty, error) {
	var res emptypb.Empty

	err := impl.bl.MarkOrderAsPaid(ctx, model.UUID(req.Id))
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (impl CheckoutServiceImplementation) CancelOrder(ctx context.Context, req *desc.CancelOrderRequest) (*emptypb.Empty, error) {
	var res emptypb.Empty

	err := impl.bl.CancelOrder(ctx, model.UUID(req.Id))
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (impl CheckoutServiceImplementation) bindModelToDescOrder(o *model.Order, is []*model.Item) *desc.Order {
	if o == nil {
		return nil
	}

	var status desc.Status
	switch o.Status {
	case model.StatusCreated:
		status = desc.Status_CREATED
	case model.StatusPaid:
		status = desc.Status_PAID
	case model.StatusFailed:
		status = desc.Status_FAILED
	case model.StatusCancelled:
		status = desc.Status_CANCELLED
	}

	var items []*desc.Item
	for _, i := range is {
		items = append(items, &desc.Item{
			Id:       string(i.Id),
			OrderId:  string(i.OrderId),
			StockId:  string(i.StockId),
			TicketId: string(i.TicketId),
			Count:    int32(i.Count),
		})
	}

	return &desc.Order{
		Id:     string(o.Id),
		UserId: string(o.UserId),
		Status: status,
		Items:  items,
	}
}

func (impl CheckoutServiceImplementation) bindModelToDescCart(cart *model.Cart) *desc.Cart {
	return &desc.Cart{
		Id:       string(cart.Id),
		UserId:   string(cart.UserId),
		TicketId: string(cart.TicketId),
		Count:    int32(cart.Count),
	}
}
