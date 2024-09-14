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
	"ticket-system/booking/internal/config"
	handlerv1 "ticket-system/booking/internal/handler/v1"
	repository "ticket-system/booking/internal/repository/postgres"
	"ticket-system/booking/internal/service"
	v1 "ticket-system/booking/pkg/v1"
	"ticket-system/lib/interceptor"
	"ticket-system/lib/logger"
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

	connPool, err := pgxpool.New(ctx, config.Data.DatabaseUrl)
	if err != nil {
		logger.Fatal("unable to connect to database", zap.Error(err))
	}
	defer connPool.Close()
	db := postgres.New(connPool)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Data.Port))
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	bls := makeBusinessLogicService(db)

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.LoggingInterceptor,
		),
	)
	v1.RegisterBookingServiceServer(s, handlerv1.NewBookingServiceImplementation(bls))
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

func makeBusinessLogicService(db *postgres.TransactionManager) service.BusinessLogic {
	sr := repository.NewStockRepository(db)
	br := repository.NewBookingRepository(db)

	ss := service.NewStockService(sr)
	bs := service.NewBookingService(br)

	bl := service.New(db, ss, bs)

	return bl
}
