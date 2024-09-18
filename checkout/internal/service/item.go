package service

import (
	"context"
	"ticket-system/checkout/internal/model"
	"ticket-system/checkout/internal/repository"
)

type itemService struct {
	itemRepository repository.ItemRepository
}

func NewItemService(eventRepository repository.ItemRepository) ItemService {
	return &itemService{
		eventRepository,
	}
}

func (s itemService) GetItem(ctx context.Context, id model.UUID) (*model.Item, error) {
	return s.itemRepository.Find(ctx, id)
}

func (s itemService) CreateItem(ctx context.Context, i *model.Item) error {
	return s.itemRepository.Create(ctx, i)
}

func (s itemService) ListItems(ctx context.Context, orderId model.UUID) ([]*model.Item, error) {
	filter := repository.FindItemsByParams{
		OrderId: orderId,
	}
	return s.itemRepository.FindBy(ctx, filter)
}
