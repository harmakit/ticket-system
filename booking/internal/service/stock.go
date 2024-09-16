package service

import (
	"context"
	"ticket-system/booking/internal/model"
	"ticket-system/booking/internal/repository"
)

type stockService struct {
	stockRepository repository.StockRepository
}

func NewStockService(stockRepository repository.StockRepository) StockService {
	return &stockService{
		stockRepository,
	}
}

func (s stockService) GetStock(ctx context.Context, uuid model.UUID) (*model.Stock, error) {
	return s.stockRepository.Find(ctx, uuid)
}

func (s stockService) ListStocks(ctx context.Context, eventId model.UUID, ticketId *model.UUID) ([]*model.Stock, error) {
	filter := repository.FindStocksByParams{
		EventId: eventId,
	}
	if ticketId != nil {
		filter.TicketId.Valid = true
		filter.TicketId.Value = *ticketId
	}

	return s.stockRepository.FindBy(ctx, filter)
}

func (s stockService) CreateStock(ctx context.Context, stock *model.Stock) error {
	return s.stockRepository.Create(ctx, stock)
}

func (s stockService) UpdateStock(ctx context.Context, stock *model.Stock) error {
	return s.stockRepository.Update(ctx, stock)
}

func (s stockService) ModifyBookedSeats(ctx context.Context, stock *model.Stock, quantity int) error {
	err := s.stockRepository.ModifyBookedSeats(ctx, stock, quantity)
	if err != nil {
		return err
	}

	if stock.SeatsBooked > stock.SeatsTotal {
		return ErrStockIsNotEnough
	}

	if stock.SeatsBooked < 0 {
		return ErrStockIsNegative
	}

	return nil
}
