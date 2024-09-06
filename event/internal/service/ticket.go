package service

import (
	"context"
	"ticket-system/event/internal/model"
	"ticket-system/event/internal/repository"
)

type ticketService struct {
	ticketRepository repository.TicketRepository
}

type TicketService interface {
	GetTicketsForEvent(ctx context.Context, event *model.Event) ([]*model.Ticket, error)
}

func NewTicketService(ticketRepository repository.TicketRepository) TicketService {
	return &ticketService{
		ticketRepository,
	}
}

func (s *ticketService) GetTicketsForEvent(ctx context.Context, event *model.Event) ([]*model.Ticket, error) {
	return s.ticketRepository.FindByEventId(ctx, string(event.Id))
}
