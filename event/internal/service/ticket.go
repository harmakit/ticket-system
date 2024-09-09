package service

import (
	"context"
	"github.com/pkg/errors"
	"ticket-system/event/internal/model"
	"ticket-system/event/internal/repository"
	workerpool "ticket-system/lib/worker-pool"
)

type ticketService struct {
	ticketRepository repository.TicketRepository
}

type TicketService interface {
	GetTicketsForEvent(ctx context.Context, event *model.Event) ([]*model.Ticket, error)
	GetTicketsForEvents(ctx context.Context, events []*model.Event) ([][]*model.Ticket, error)
}

func NewTicketService(ticketRepository repository.TicketRepository) TicketService {
	return &ticketService{
		ticketRepository,
	}
}

func (s *ticketService) GetTicketsForEvent(ctx context.Context, event *model.Event) ([]*model.Ticket, error) {
	tickets, err := s.ticketRepository.FindByEventId(ctx, string(event.Id))
	if err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			return nil, nil
		}
	}

	return tickets, err
}

type GetEventTicketsTask struct {
	ctx   context.Context
	event *model.Event
	idx   int
}

type GetEventTicketsTaskResult struct {
	tickets []*model.Ticket
	idx     int
}

func (s *ticketService) GetTicketsForEvents(ctx context.Context, events []*model.Event) ([][]*model.Ticket, error) {
	const workersCount = 3
	wp, wpResults := workerpool.New[GetEventTicketsTask, GetEventTicketsTaskResult](ctx, workersCount)

	wpTasks := make([]workerpool.Task[GetEventTicketsTask, GetEventTicketsTaskResult], len(events))

	for i, event := range events {
		wpTasks[i] = workerpool.Task[GetEventTicketsTask, GetEventTicketsTaskResult]{
			Args: GetEventTicketsTask{
				ctx:   ctx,
				event: event,
				idx:   i,
			},
			Callback: func(task GetEventTicketsTask) (GetEventTicketsTaskResult, error) {
				result := GetEventTicketsTaskResult{
					idx: task.idx,
				}
				tickets, err := s.GetTicketsForEvent(task.ctx, task.event)
				if err == nil {
					result.tickets = tickets
				}
				return result, err
			},
		}
	}

	err := wp.Submit(ctx, wpTasks)
	if err != nil {
		return nil, err
	}

	tickets := make([][]*model.Ticket, len(events))
	for i := 0; i < len(events); i++ {
		task := <-wpResults
		if task.Err != nil {
			return nil, task.Err
		}
		tickets[task.Result.idx] = task.Result.tickets
	}

	return tickets, nil
}
