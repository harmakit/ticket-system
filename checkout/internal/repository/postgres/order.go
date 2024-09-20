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

type orderRepository struct {
	transactionManager *postgres.TransactionManager
}

var _ repository.OrderRepository = (*orderRepository)(nil)

func NewOrderRepository(transactionManager *postgres.TransactionManager) repository.OrderRepository {
	return &orderRepository{transactionManager}
}

func (r orderRepository) getQueryBuilder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}

func (r orderRepository) bindSchemaToModel(e *schema.Order) *model.Order {
	order := &model.Order{
		Id:     model.UUID(e.Id),
		UserId: model.UUID(e.UserId),
		Status: model.OrderStatus(e.Status),
	}

	return order
}

func (r orderRepository) Find(ctx context.Context, id model.UUID) (*model.Order, error) {
	db := r.transactionManager.GetQueryEngine(ctx)

	query := r.getQueryBuilder().Select(schema.OrderColumns...).From(schema.OrderTable).Where(sq.Eq{"id": id})

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, _ := db.Query(ctx, rawQuery, args...)
	order, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[schema.Order])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNoRows
		}
		return nil, err
	}

	return r.bindSchemaToModel(&order), err
}

func (r orderRepository) FindBy(ctx context.Context, params repository.FindOrdersByParams) ([]*model.Order, error) {
	db := r.transactionManager.GetQueryEngine(ctx)

	query := r.getQueryBuilder().Select(schema.OrderColumns...).From(schema.OrderTable)
	if params.UserId.Valid {
		query = query.Where(sq.Eq{"user_id": params.UserId.Value})
	}
	if params.Limit > 0 {
		query = query.Limit(uint64(params.Limit))
	}
	if params.Offset > 0 {
		query = query.Offset(uint64(params.Offset))
	}

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, _ := db.Query(ctx, rawQuery, args...)
	orders, err := pgx.CollectRows(rows, pgx.RowToStructByName[schema.Order])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNoRows
		}
		return nil, err
	}

	res := make([]*model.Order, len(orders))
	for i, event := range orders {
		res[i] = r.bindSchemaToModel(&event)
	}

	return res, err
}

func (r orderRepository) Create(ctx context.Context, o *model.Order) error {
	db := r.transactionManager.GetQueryEngine(ctx)

	query := r.getQueryBuilder().Insert(schema.OrderTable).Columns(schema.OrderColumns...).
		Values(sq.Expr(repository.NewUUID), o.UserId, o.Status).
		Suffix("RETURNING *")

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	rows, _ := db.Query(ctx, rawQuery, args...)
	no, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[schema.Order])
	if err != nil {
		return err
	}

	*o = *r.bindSchemaToModel(&no)

	return nil
}

func (r orderRepository) Update(ctx context.Context, o *model.Order) error {
	db := r.transactionManager.GetQueryEngine(ctx)

	query := r.getQueryBuilder().Update(schema.OrderTable).
		SetMap(map[string]interface{}{
			"status": o.Status,
		}).
		Where(sq.Eq{"id": o.Id}).
		Suffix("RETURNING *")

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	rows, _ := db.Query(ctx, rawQuery, args...)

	no, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[schema.Order])
	if err != nil {
		return err
	}

	*o = *r.bindSchemaToModel(&no)

	return nil
}
