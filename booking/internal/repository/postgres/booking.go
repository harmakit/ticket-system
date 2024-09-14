package postgres

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"ticket-system/booking/internal/model"
	"ticket-system/booking/internal/repository"
	"ticket-system/booking/internal/repository/schema"
	"ticket-system/lib/query-engine/postgres"
	"time"
)

type bookingRepository struct {
	transactionManager *postgres.TransactionManager
}

var _ repository.BookingRepository = (*bookingRepository)(nil)

func NewBookingRepository(transactionManager *postgres.TransactionManager) repository.BookingRepository {
	return &bookingRepository{transactionManager}
}

func (r bookingRepository) getQueryBuilder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}

func (r bookingRepository) bindSchemaToModel(b *schema.Booking) *model.Booking {
	return &model.Booking{
		Id:        model.UUID(b.Id),
		StockId:   model.UUID(b.StockId),
		UserId:    model.UUID(b.UserId),
		OrderId:   model.UUID(b.OrderId),
		Count:     b.Count,
		CreatedAt: b.CreatedAt,
		ExpiredAt: b.ExpiredAt,
	}
}

func (r bookingRepository) Find(ctx context.Context, id model.UUID) (*model.Booking, error) {
	db := r.transactionManager.GetQueryEngine(ctx)

	query := r.getQueryBuilder().Select(schema.BookingColumns...).From(schema.BookingTable).Where(sq.Eq{"id": id})

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, _ := db.Query(ctx, rawQuery, args...)
	booking, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[schema.Booking])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNoRows
		}
		return nil, err
	}

	return r.bindSchemaToModel(&booking), err
}

func (r bookingRepository) FindBy(ctx context.Context, params repository.FindBookingsByParams) ([]*model.Booking, error) {
	db := r.transactionManager.GetQueryEngine(ctx)

	query := r.getQueryBuilder().Select(schema.BookingColumns...).From(schema.BookingTable)
	query = query.Where(sq.Eq{"stock_id": params.StockId})
	query = query.Where(sq.Eq{"user_id": params.UserId})
	query = query.Where(sq.Eq{"order_id": params.OrderId})

	if !params.WithExpired {
		query = query.Where(sq.Gt{"expired_at": time.Now()})
	}

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, _ := db.Query(ctx, rawQuery, args...)
	bookings, err := pgx.CollectRows(rows, pgx.RowToStructByName[schema.Booking])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNoRows
		}
		return nil, err
	}

	res := make([]*model.Booking, len(bookings))
	for i, booking := range bookings {
		res[i] = r.bindSchemaToModel(&booking)
	}

	return res, err
}

func (r bookingRepository) Create(ctx context.Context, b *model.Booking) error {
	db := r.transactionManager.GetQueryEngine(ctx)

	query := r.getQueryBuilder().Insert(schema.BookingTable).Columns(schema.BookingColumns...).
		Values(sq.Expr("gen_random_uuid()"), b.StockId, b.UserId, b.OrderId, b.Count, b.CreatedAt, b.ExpiredAt).
		Suffix("RETURNING *")

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	rows, _ := db.Query(ctx, rawQuery, args...)
	nb, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[schema.Booking])
	if err != nil {
		return err
	}

	*b = *r.bindSchemaToModel(&nb)

	return nil
}