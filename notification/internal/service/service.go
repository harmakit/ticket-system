package service

import (
	"ticket-system/notification/internal/model"
)

type LoggerService interface {
	LogOrder(order *model.Order)
}

type BusinessLogic struct {
	loggerService LoggerService
}

func New(loggerService LoggerService) BusinessLogic {
	return BusinessLogic{
		loggerService,
	}
}

func (s *BusinessLogic) NotifyOrderChanges(order *model.Order) error {
	s.loggerService.LogOrder(order)
	// ... notify via email or whatever else
	return nil
}
