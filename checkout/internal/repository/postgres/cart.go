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

type cartRepository struct {
	transactionManager *postgres.TransactionManager
}

var _ repository.CartRepository = (*cartRepository)(nil)

func NewCartRepository(transactionManager *postgres.TransactionManager) repository.CartRepository {
	return &cartRepository{transactionManager}
}

func (r cartRepository) getQueryBuilder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}

func (r cartRepository) bindSchemaToModel(e *schema.Cart) *model.Cart {
	cart := &model.Cart{
		Id:       model.UUID(e.Id),
		UserId:   model.UUID(e.UserId),
		TicketId: model.UUID(e.TicketId),
		Count:    e.Count,
	}

	return cart
}

func (r cartRepository) Find(ctx context.Context, id model.UUID) (*model.Cart, error) {
	db := r.transactionManager.GetQueryEngine(ctx)

	query := r.getQueryBuilder().Select(schema.CartColumns...).From(schema.CartTable).Where(sq.Eq{"id": id})

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, _ := db.Query(ctx, rawQuery, args...)
	cart, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[schema.Cart])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNoRows
		}
		return nil, err
	}

	return r.bindSchemaToModel(&cart), err
}

func (r cartRepository) FindBy(ctx context.Context, filter repository.FindCartsByParams) ([]*model.Cart, error) {
	db := r.transactionManager.GetQueryEngine(ctx)

	query := r.getQueryBuilder().Select(schema.CartColumns...).From(schema.CartTable)
	if filter.UserId.Valid {
		query = query.Where(sq.Eq{"user_id": filter.UserId.Value})
	}
	if filter.TicketId.Valid {
		query = query.Where(sq.Eq{"ticket_id": filter.TicketId.Value})
	}

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, _ := db.Query(ctx, rawQuery, args...)
	carts, err := pgx.CollectRows(rows, pgx.RowToStructByName[schema.Cart])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNoRows
		}
		return nil, err
	}

	res := make([]*model.Cart, len(carts))
	for i, cart := range carts {
		res[i] = r.bindSchemaToModel(&cart)
	}

	return res, err
}

func (r cartRepository) Create(ctx context.Context, c *model.Cart) error {
	db := r.transactionManager.GetQueryEngine(ctx)

	query := r.getQueryBuilder().Insert(schema.CartTable).Columns(schema.CartColumns...).
		Values(sq.Expr(repository.NewUUID), c.UserId, c.TicketId, c.Count).
		Suffix("RETURNING *")

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	rows, _ := db.Query(ctx, rawQuery, args...)
	nc, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[schema.Cart])
	if err != nil {
		return err
	}

	*c = *r.bindSchemaToModel(&nc)

	return nil
}

func (r cartRepository) Update(ctx context.Context, c *model.Cart) error {
	db := r.transactionManager.GetQueryEngine(ctx)

	query := r.getQueryBuilder().Update(schema.CartTable).
		SetMap(map[string]interface{}{
			"count": c.Count,
		}).
		Where(sq.Eq{"id": c.Id}).
		Suffix("RETURNING *")

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	rows, _ := db.Query(ctx, rawQuery, args...)

	nc, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[schema.Cart])
	if err != nil {
		return err
	}

	*c = *r.bindSchemaToModel(&nc)

	return nil
}

func (r cartRepository) Delete(ctx context.Context, id model.UUID) error {
	db := r.transactionManager.GetQueryEngine(ctx)

	query := r.getQueryBuilder().Delete(schema.CartTable).Where(sq.Eq{"id": id})

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, rawQuery, args...)

	return err
}
