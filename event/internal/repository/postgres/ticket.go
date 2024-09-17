package postgres

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"ticket-system/event/internal/model"
	"ticket-system/event/internal/repository"
	"ticket-system/event/internal/repository/schema"
	"ticket-system/lib/query-engine/postgres"
)

type ticketRepository struct {
	transactionManager *postgres.TransactionManager
}

var _ repository.LocationRepository = (*locationRepository)(nil)

func NewTicketRepository(transactionManager *postgres.TransactionManager) repository.TicketRepository {
	return &ticketRepository{transactionManager}
}

func (r ticketRepository) getQueryBuilder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}

func (r ticketRepository) bindSchemaToModel(t *schema.Ticket) *model.Ticket {
	return &model.Ticket{
		Id:      model.UUID(t.Id),
		EventId: model.UUID(t.EventId),
		Name:    t.Name,
		Price:   t.Price,
	}
}

func (r ticketRepository) Find(ctx context.Context, id model.UUID) (*model.Ticket, error) {
	db := r.transactionManager.GetQueryEngine(ctx)

	query := r.getQueryBuilder().Select(schema.TicketColumns...).From(schema.TicketTable).Where(sq.Eq{"id": id})

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, _ := db.Query(ctx, rawQuery, args...)
	ticket, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[schema.Ticket])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNoRows
		}
		return nil, err
	}

	return r.bindSchemaToModel(&ticket), err
}

func (r ticketRepository) FindByEventId(ctx context.Context, eventId model.UUID) ([]*model.Ticket, error) {
	db := r.transactionManager.GetQueryEngine(ctx)

	query := r.getQueryBuilder().Select(schema.TicketColumns...).From(schema.TicketTable).Where(sq.Eq{"event_id": eventId})

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, _ := db.Query(ctx, rawQuery, args...)
	tickets, err := pgx.CollectRows(rows, pgx.RowToStructByName[schema.Ticket])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNoRows
		}
		return nil, err
	}

	res := make([]*model.Ticket, len(tickets))
	for i, ticket := range tickets {
		res[i] = r.bindSchemaToModel(&ticket)
	}

	return res, err
}

func (r ticketRepository) Create(ctx context.Context, t *model.Ticket) error {
	db := r.transactionManager.GetQueryEngine(ctx)

	query := r.getQueryBuilder().Insert(schema.TicketTable).Columns(schema.TicketColumns...).
		Values(sq.Expr(repository.NewUUID), t.EventId, t.Name, t.Price).
		Suffix("RETURNING *")

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	rows, _ := db.Query(ctx, rawQuery, args...)
	nt, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[schema.Ticket])
	if err != nil {
		return err
	}

	*t = *r.bindSchemaToModel(&nt)

	return nil
}
