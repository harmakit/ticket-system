package service

import (
	"fmt"
	"ticket-system/lib/logger"
	"ticket-system/notification/internal/model"
)

type loggerService struct {
}

func NewLoggerService() LoggerService {
	return &loggerService{}
}

func (s loggerService) LogOrder(order *model.Order) {
	message := fmt.Sprintf("order (%s) updated! status: %s", order.Id, order.Status)
	logger.Info(message)
}
