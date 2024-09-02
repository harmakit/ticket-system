package postgres

import (
	"context"
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

func (r eventRepository) bindSchemaToModel(event *schema.Event) *model.Event {
	return &model.Event{
		Id:          model.UUID(event.Id),
		Date:        event.Date,
		Duration:    event.Duration,
		Name:        event.Name,
		Description: event.Description,
		LocationId:  event.LocationId,
	}
}

func (r eventRepository) Find(ctx context.Context, id string) (*model.Event, error) {
	db := r.transactionManager.GetQueryEngine(ctx)

	query := r.getQueryBuilder().Select(schema.EventColumns...).From(schema.EventTable).Where(sq.Eq{"id": id})

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, _ := db.Query(ctx, rawQuery, args)
	event, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[schema.Event])
	if err != nil {
		return nil, err
	}

	return r.bindSchemaToModel(&event), err
}