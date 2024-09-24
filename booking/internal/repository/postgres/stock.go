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
)

type stockRepository struct {
	transactionManager *postgres.TransactionManager
}

var _ repository.StockRepository = (*stockRepository)(nil)

func NewStockRepository(transactionManager *postgres.TransactionManager) repository.StockRepository {
	return &stockRepository{transactionManager}
}

func (r stockRepository) getQueryBuilder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}

func (r stockRepository) bindSchemaToModel(s *schema.Stock) *model.Stock {
	return &model.Stock{
		Id:          model.UUID(s.Id),
		EventId:     model.UUID(s.EventId),
		TicketId:    model.UUID(s.TicketId),
		SeatsTotal:  s.SeatsTotal,
		SeatsBooked: s.SeatsBooked,
	}
}

func (r stockRepository) Find(ctx context.Context, id model.UUID) (*model.Stock, error) {
	db := r.transactionManager.GetQueryEngine(ctx)

	query := r.getQueryBuilder().Select(schema.StockColumns...).From(schema.StockTable).Where(sq.Eq{"id": id})

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, _ := db.Query(ctx, rawQuery, args...)
	stock, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[schema.Stock])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNoRows
		}
		return nil, err
	}

	return r.bindSchemaToModel(&stock), err
}

func (r stockRepository) Create(ctx context.Context, s *model.Stock) error {
	db := r.transactionManager.GetQueryEngine(ctx)

	query := r.getQueryBuilder().Insert(schema.StockTable).Columns(schema.StockColumns...).
		Values(sq.Expr(NewUUID), s.EventId, s.TicketId, s.SeatsTotal, s.SeatsBooked).
		Suffix("RETURNING *")

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	rows, _ := db.Query(ctx, rawQuery, args...)
	ns, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[schema.Stock])
	if err != nil {
		return err
	}

	*s = *r.bindSchemaToModel(&ns)

	return nil
}

func (r stockRepository) FindBy(ctx context.Context, filter repository.FindStocksByParams) ([]*model.Stock, error) {
	db := r.transactionManager.GetQueryEngine(ctx)

	query := r.getQueryBuilder().Select(schema.StockColumns...).From(schema.StockTable)
	query = query.Where(sq.Eq{"event_id": filter.EventId})
	if filter.TicketId.Valid {
		query = query.Where(sq.Eq{"ticket_id": repository.NullUUID(filter.TicketId)})
	}

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, _ := db.Query(ctx, rawQuery, args...)
	stocks, err := pgx.CollectRows(rows, pgx.RowToStructByName[schema.Stock])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNoRows
		}
		return nil, err
	}

	res := make([]*model.Stock, len(stocks))
	for i, stock := range stocks {
		res[i] = r.bindSchemaToModel(&stock)
	}

	return res, err
}

func (r stockRepository) Update(ctx context.Context, s *model.Stock) error {
	db := r.transactionManager.GetQueryEngine(ctx)

	query := r.getQueryBuilder().Update(schema.StockTable).
		SetMap(map[string]interface{}{
			"seats_booked": s.SeatsBooked,
		}).
		Where(sq.Eq{"id": s.Id}).
		Suffix("RETURNING *")

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	rows, _ := db.Query(ctx, rawQuery, args...)

	ns, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[schema.Stock])
	if err != nil {
		return err
	}

	*s = *r.bindSchemaToModel(&ns)

	return nil
}

func (r stockRepository) ModifyBookedSeats(ctx context.Context, s *model.Stock, quantity int) error {
	db := r.transactionManager.GetQueryEngine(ctx)

	query := r.getQueryBuilder().Update(schema.StockTable).
		SetMap(map[string]interface{}{
			"seats_booked": sq.Expr("seats_booked + ?", quantity),
		}).
		Where(sq.Eq{"id": s.Id}).
		Suffix("RETURNING *")

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	rows, _ := db.Query(ctx, rawQuery, args...)
	ns, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[schema.Stock])
	if err != nil {
		return err
	}

	*s = *r.bindSchemaToModel(&ns)

	return nil
}

func (r stockRepository) Delete(ctx context.Context, id model.UUID) error {
	db := r.transactionManager.GetQueryEngine(ctx)

	query := r.getQueryBuilder().Delete(schema.StockTable).Where(sq.Eq{"id": id})

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, rawQuery, args...)

	return err
}
