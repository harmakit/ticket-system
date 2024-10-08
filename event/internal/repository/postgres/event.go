package postgres

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"ticket-system/event/internal/model"
	"ticket-system/event/internal/repository"
	"ticket-system/event/internal/repository/schema"
	"ticket-system/lib/query-engine/postgres"
)

type eventRepository struct {
	transactionManager *postgres.TransactionManager
}

var _ repository.EventRepository = (*eventRepository)(nil)

func NewEventRepository(transactionManager *postgres.TransactionManager) repository.EventRepository {
	return &eventRepository{transactionManager}
}

func (r eventRepository) getQueryBuilder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}

func (r eventRepository) bindSchemaToModel(e *schema.Event) *model.Event {
	event := &model.Event{
		Id:       model.UUID(e.Id),
		Date:     e.StartDate,
		Duration: e.Duration,
		Name:     e.Name,
	}
	if e.Description.Valid {
		event.Description = model.NullString{
			Value: e.Description.String,
			Valid: true,
		}
	}
	if e.LocationId.Valid {
		event.LocationId = model.NullUUID{
			Value: model.UUID(e.LocationId.String),
			Valid: true,
		}
	}

	return event
}

func (r eventRepository) Find(ctx context.Context, id model.UUID) (*model.Event, error) {
	db := r.transactionManager.GetQueryEngine(ctx)

	query := r.getQueryBuilder().Select(schema.EventColumns...).From(schema.EventTable).Where(sq.Eq{"id": id})

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, _ := db.Query(ctx, rawQuery, args...)
	event, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[schema.Event])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNoRows
		}
		return nil, err
	}

	return r.bindSchemaToModel(&event), err
}

func (r eventRepository) FindBy(ctx context.Context, params repository.FindEventsByParams) ([]*model.Event, error) {
	db := r.transactionManager.GetQueryEngine(ctx)

	query := r.getQueryBuilder().Select(schema.EventColumns...).From(schema.EventTable)
	if params.LocationId.Valid {
		query = query.Where(sq.Eq{"location_id": repository.NullUUID(params.LocationId)})
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
	events, err := pgx.CollectRows(rows, pgx.RowToStructByName[schema.Event])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNoRows
		}
		return nil, err
	}

	res := make([]*model.Event, len(events))
	for i, event := range events {
		res[i] = r.bindSchemaToModel(&event)
	}

	return res, err
}

func (r eventRepository) Create(ctx context.Context, e *model.Event) error {
	db := r.transactionManager.GetQueryEngine(ctx)

	query := r.getQueryBuilder().Insert(schema.EventTable).Columns(schema.EventColumns...).
		Values(sq.Expr(NewUUID), e.Date, e.Duration, e.Name, repository.NullString(e.Description), repository.NullUUID(e.LocationId)).
		Suffix("RETURNING *")

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	rows, _ := db.Query(ctx, rawQuery, args...)
	ne, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[schema.Event])
	if err != nil {
		return err
	}

	*e = *r.bindSchemaToModel(&ne)

	return nil
}
