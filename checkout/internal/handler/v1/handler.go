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

func (impl CheckoutServiceImplementation) PlaceOrder(ctx context.Context, req *desc.PlaceOrderRequest) (*emptypb.Empty, error) {
	var res emptypb.Empty

	err := impl.bl.PlaceOrder(ctx, model.UUID(req.UserId))
	if err != nil {
		return nil, err
	}

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
			Id:      string(i.Id),
			OrderId: string(i.OrderId),
			StockId: string(i.StockId),
			Count:   int32(i.Count),
		})
	}

	return &desc.Order{
		Id:     string(o.Id),
		UserId: string(o.UserId),
		Status: status,
		Items:  items,
	}
}
