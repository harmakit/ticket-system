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
	"ticket-system/event/internal/client/booking"
	"ticket-system/event/internal/config"
	handlerv1 "ticket-system/event/internal/handler/v1"
	repository "ticket-system/event/internal/repository/postgres"
	"ticket-system/event/internal/service"
	v1 "ticket-system/event/pkg/v1"
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
	v1.RegisterEventServiceServer(s, handlerv1.NewEventServiceImplementation(bls))
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
	er := repository.NewEventRepository(db)
	lr := repository.NewLocationRepository(db)
	tr := repository.NewTicketRepository(db)

	es := service.NewEventService(er)
	ls := service.NewLocationService(lr)
	ts := service.NewTicketService(tr)

	bs := booking.NewService()

	bl := service.New(db, es, ls, ts, bs)

	return bl
}
