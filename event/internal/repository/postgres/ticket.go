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

func (r ticketRepository) bindSchemaToModel(ticket *schema.Ticket) *model.Ticket {
	return &model.Ticket{
		EventId: model.UUID(ticket.EventId),
		Name:    ticket.Name,
		Price:   ticket.Price,
	}
}

func (r ticketRepository) Find(ctx context.Context, id string) (*model.Ticket, error) {
	db := r.transactionManager.GetQueryEngine(ctx)

	query := r.getQueryBuilder().Select(schema.TicketColumns...).From(schema.TicketTable).Where(sq.Eq{"id": id})

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, _ := db.Query(ctx, rawQuery, args)
	Ticket, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[schema.Ticket])
	if err != nil {
		return nil, err
	}

	return r.bindSchemaToModel(&Ticket), err
}
