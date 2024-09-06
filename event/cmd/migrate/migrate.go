package main

import "os/signal"

import (
	"context"
	"database/sql"
	gomigrate "github.com/golang-migrate/migrate/v4"
	gomigartedriver "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"os"
	"ticket-system/event/internal/config"
	"ticket-system/lib/logger"
)

func main() {
	var err error

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err = config.Init()
	if err != nil {
		panic("unable to initialize config: " + err.Error())
	}

	err = logger.Init(config.Data.IsDev())
	if err != nil {
		panic("unable to initialize logger: " + err.Error())
	}

	connPool, err := pgxpool.New(ctx, config.Data.DatabaseUrl)
	if err != nil {
		logger.Fatal("unable to connect to database", zap.Error(err))
	}
	defer connPool.Close()

	db, err := sql.Open("pgx/v5", config.Data.DatabaseUrl)
	if err != nil {
		logger.Fatal("unable to connect to database", zap.Error(err))
	}

	driver, err := gomigartedriver.WithInstance(db, &gomigartedriver.Config{})
	if err != nil {
		logger.Fatal("unable to initialize driver", zap.Error(err))
	}

	m, err := gomigrate.NewWithDatabaseInstance(
		"file:///migrations",
		"postgres", driver)
	if err != nil {
		logger.Fatal("unable to initialize migrator", zap.Error(err))
	}

	err = m.Up()
	if err != nil {
		logger.Fatal("unable to run migrations", zap.Error(err))
	}

	logger.Info("database migrated")
}
