package postgres

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"ticket-system/checkout/internal/model"
	"ticket-system/checkout/internal/repository"
	"ticket-system/checkout/internal/repository/schema"
	"ticket-system/lib/query-engine/postgres"
)

type itemRepository struct {
	transactionManager *postgres.TransactionManager
}

var _ repository.ItemRepository = (*itemRepository)(nil)

func NewItemRepository(transactionManager *postgres.TransactionManager) repository.ItemRepository {
	return &itemRepository{transactionManager}
}

func (r itemRepository) getQueryBuilder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}

func (r itemRepository) bindSchemaToModel(e *schema.Item) *model.Item {
	item := &model.Item{
		Id:       model.UUID(e.Id),
		OrderId:  model.UUID(e.OrderId),
		StockId:  model.UUID(e.StockId),
		TicketId: model.UUID(e.TicketId),
		Count:    e.Count,
	}

	return item
}

func (r itemRepository) Find(ctx context.Context, id model.UUID) (*model.Item, error) {
	db := r.transactionManager.GetQueryEngine(ctx)

	query := r.getQueryBuilder().Select(schema.ItemColumns...).From(schema.ItemTable).Where(sq.Eq{"id": id})

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, _ := db.Query(ctx, rawQuery, args...)
	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[schema.Item])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNoRows
		}
		return nil, err
	}

	return r.bindSchemaToModel(&item), err
}

func (r itemRepository) FindBy(ctx context.Context, filter repository.FindItemsByParams) ([]*model.Item, error) {
	db := r.transactionManager.GetQueryEngine(ctx)

	query := r.getQueryBuilder().Select(schema.ItemColumns...).From(schema.ItemTable)
	query = query.Where(sq.Eq{"order_id": filter.OrderId})

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, _ := db.Query(ctx, rawQuery, args...)
	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[schema.Item])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNoRows
		}
		return nil, err
	}

	res := make([]*model.Item, len(items))
	for i, item := range items {
		res[i] = r.bindSchemaToModel(&item)
	}

	return res, err
}

func (r itemRepository) Create(ctx context.Context, i *model.Item) error {
	db := r.transactionManager.GetQueryEngine(ctx)

	query := r.getQueryBuilder().Insert(schema.ItemTable).Columns(schema.ItemColumns...).
		Values(sq.Expr(NewUUID), i.OrderId, i.StockId, i.TicketId, i.Count).
		Suffix("RETURNING *")

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	rows, _ := db.Query(ctx, rawQuery, args...)
	ni, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[schema.Item])
	if err != nil {
		return err
	}

	*i = *r.bindSchemaToModel(&ni)

	return nil
}
