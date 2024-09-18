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
	"ticket-system/checkout/internal/client/booking"
	"ticket-system/checkout/internal/client/event"
	"ticket-system/checkout/internal/config"
	handlerv1 "ticket-system/checkout/internal/handler/v1"
	repository "ticket-system/checkout/internal/repository/postgres"
	"ticket-system/checkout/internal/service"
	v1 "ticket-system/checkout/pkg/v1"
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
	v1.RegisterCheckoutServiceServer(s, handlerv1.NewCheckoutServiceImplementation(bls))
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
	or := repository.NewOrderRepository(db)
	ir := repository.NewItemRepository(db)
	cr := repository.NewCartRepository(db)

	ors := service.NewOrderService(or)
	is := service.NewItemService(ir)
	cs := service.NewCartService(cr)

	bs := booking.NewService()
	es := event.NewService()

	bl := service.New(db, ors, is, cs, bs, es)

	return bl
}
