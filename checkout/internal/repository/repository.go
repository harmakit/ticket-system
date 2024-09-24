package repository

import (
	"context"
	"database/sql"
	"ticket-system/checkout/internal/model"
)

type OrderRepository interface {
	Find(ctx context.Context, id model.UUID) (*model.Order, error)
	Create(ctx context.Context, o *model.Order) error
	Update(ctx context.Context, order *model.Order) error
	FindBy(ctx context.Context, filter FindOrdersByParams) ([]*model.Order, error)
}

type ItemRepository interface {
	Find(ctx context.Context, id model.UUID) (*model.Item, error)
	FindBy(ctx context.Context, filter FindItemsByParams) ([]*model.Item, error)
	Create(ctx context.Context, o *model.Item) error
}

type CartRepository interface {
	Find(ctx context.Context, id model.UUID) (*model.Cart, error)
	FindBy(ctx context.Context, filter FindCartsByParams) ([]*model.Cart, error)
	Create(ctx context.Context, c *model.Cart) error
	Update(ctx context.Context, c *model.Cart) error
	Delete(ctx context.Context, id model.UUID) error
}

func NullString(s model.NullString) sql.NullString {
	return sql.NullString{
		String: s.Value,
		Valid:  s.Valid,
	}
}

func NullUUID(s model.NullUUID) sql.NullString {
	return sql.NullString{
		String: string(s.Value),
		Valid:  s.Valid,
	}
}

type FindItemsByParams struct {
	OrderId model.UUID
}

type FindCartsByParams struct {
	UserId   model.NullUUID
	TicketId model.NullUUID
}

type FindOrdersByParams struct {
	UserId model.NullUUID
	Limit  int
	Offset int
}
