package service

import (
	"context"
	"github.com/pkg/errors"
	"ticket-system/event/internal/model"
	"ticket-system/event/internal/repository"
	workerpool "ticket-system/lib/worker-pool"
)

type locationService struct {
	locationRepository repository.LocationRepository
}

type LocationService interface {
	GetLocation(ctx context.Context, uuid model.UUID) (*model.Location, error)
	GetLocationsForEvents(ctx context.Context, events []*model.Event) ([]*model.Location, error)
}

func NewLocationService(locationRepository repository.LocationRepository) LocationService {
	return &locationService{
		locationRepository,
	}
}

func (s *locationService) GetLocation(ctx context.Context, uuid model.UUID) (*model.Location, error) {
	location, err := s.locationRepository.Find(ctx, string(uuid))
	if err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			return nil, nil
		}
	}

	return location, err
}

type GetEventLocationTask struct {
	ctx   context.Context
	event *model.Event
	idx   int
}

type GetEventLocationTaskResult struct {
	location *model.Location
	idx      int
}

func (s *locationService) GetLocationsForEvents(ctx context.Context, events []*model.Event) ([]*model.Location, error) {
	const workersCount = 3
	wp, wpResults := workerpool.New[GetEventLocationTask, GetEventLocationTaskResult](ctx, workersCount)

	wpTasks := make([]workerpool.Task[GetEventLocationTask, GetEventLocationTaskResult], len(events))

	for i, event := range events {
		wpTasks[i] = workerpool.Task[GetEventLocationTask, GetEventLocationTaskResult]{
			Args: GetEventLocationTask{
				ctx:   ctx,
				event: event,
				idx:   i,
			},
			Callback: func(task GetEventLocationTask) (GetEventLocationTaskResult, error) {
				result := GetEventLocationTaskResult{
					idx: task.idx,
				}
				location, err := s.GetLocation(task.ctx, model.UUID(task.event.LocationId.String))
				if err == nil {
					result.location = location
				}
				return result, err
			},
		}
	}

	err := wp.Submit(ctx, wpTasks)
	if err != nil {
		return nil, err
	}

	locations := make([]*model.Location, len(events))
	for i := 0; i < len(events); i++ {
		task := <-wpResults
		if task.Err != nil {
			return nil, task.Err
		}
		locations[task.Result.idx] = task.Result.location
	}

	return locations, nil
}
