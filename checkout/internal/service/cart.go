package service

import (
	"context"
	"ticket-system/checkout/internal/model"
	"ticket-system/checkout/internal/repository"
)

type cartService struct {
	cartRepository repository.CartRepository
}

func NewCartService(eventRepository repository.CartRepository) CartService {
	return &cartService{
		eventRepository,
	}
}

func (s cartService) GetCart(ctx context.Context, id model.UUID) (*model.Cart, error) {
	return s.cartRepository.Find(ctx, id)
}

func (s cartService) CreateCart(ctx context.Context, c *model.Cart) error {
	return s.cartRepository.Create(ctx, c)
}

func (s cartService) ListCarts(ctx context.Context, userId *model.UUID, ticketId *model.UUID) ([]*model.Cart, error) {
	filter := repository.FindCartsByParams{}
	if userId != nil {
		filter.UserId = model.NullUUID{Valid: true, Value: *userId}
	}
	if ticketId != nil {
		filter.TicketId = model.NullUUID{Valid: true, Value: *ticketId}
	}
	return s.cartRepository.FindBy(ctx, filter)
}

func (s cartService) UpdateCart(ctx context.Context, c *model.Cart) error {
	if c.Count <= 0 {
		return s.cartRepository.Delete(ctx, c.Id)
	}

	return s.cartRepository.Update(ctx, c)
}

func (s cartService) DeleteCart(ctx context.Context, c *model.Cart) error {
	return s.cartRepository.Delete(ctx, c.Id)
}
