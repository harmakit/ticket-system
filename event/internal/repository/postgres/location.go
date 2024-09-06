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

type locationRepository struct {
	transactionManager *postgres.TransactionManager
}

var _ repository.LocationRepository = (*locationRepository)(nil)

func NewLocationRepository(transactionManager *postgres.TransactionManager) repository.LocationRepository {
	return &locationRepository{transactionManager}
}

func (r locationRepository) getQueryBuilder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}

func (r locationRepository) bindSchemaToModel(location *schema.Location) *model.Location {
	return &model.Location{
		Id:      model.UUID(location.Id),
		Name:    location.Name,
		Address: location.Address,
		Lat:     location.Lat,
		Lng:     location.Lng,
	}
}

func (r locationRepository) Find(ctx context.Context, id string) (*model.Location, error) {
	db := r.transactionManager.GetQueryEngine(ctx)

	query := r.getQueryBuilder().Select(schema.LocationColumns...).From(schema.LocationTable).Where(sq.Eq{"id": id})

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, _ := db.Query(ctx, rawQuery, args)
	location, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[schema.Location])
	if err != nil {
		return nil, err
	}

	return r.bindSchemaToModel(&location), err
}
