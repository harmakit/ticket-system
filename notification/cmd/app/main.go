package main

import (
	"context"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"ticket-system/lib/logger"
	"ticket-system/notification/internal/config"
	handlerv1 "ticket-system/notification/internal/handler/v1"
	"ticket-system/notification/internal/message"
	"ticket-system/notification/internal/service"
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

	bls := makeBusinessLogicService()
	h := handlerv1.NewMessageHandlerImplementation(bls)

	om := message.NewOrderMessenger(h)
	runConsumer(ctx, om)

	<-ctx.Done()
	closeMessengers(om)
}

func closeMessengers(ms ...message.Messenger) {
	for _, m := range ms {
		err := m.Close()
		if err != nil {
			logger.Fatal("unable to close messenger", zap.Error(err))
		}
	}
}

func runConsumer(ctx context.Context, m message.Messenger) {
	errHandler := func(err error) {
		logger.Error("unable to consume messages", zap.Error(err))
	}
	go func() {
		m.Consume(ctx, errHandler)
	}()
}

func makeBusinessLogicService() service.BusinessLogic {
	ls := service.NewLoggerService()

	bl := service.New(ls)

	return bl
}
