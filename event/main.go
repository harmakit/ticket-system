package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"os/signal"
	"ticket-system/events/internal/config"
	"ticket-system/events/internal/logger"
	"ticket-system/lib/query-engine/postgres"
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

	connPool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Fatal("unable to connect to database", zap.Error(err))
	}
	defer connPool.Close()
	db := postgres.New(connPool)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Data.Port))
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	s := grpc.NewServer()
	reflection.Register(s)

	logger.Info("server listening", zap.String("address", lis.Addr().String()))

	go func() {
		logger.Info("starting server", zap.String("address", lis.Addr().String()))
		err = s.Serve(lis)
		if err != nil {
			logger.Fatal("failed to serve", zap.Error(err))
		}
	}()

	<-ctx.Done()
	s.GracefulStop()
}
