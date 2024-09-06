package service

import (
	"context"
	"ticket-system/event/internal/model"
	"ticket-system/event/internal/repository"
)

type locationService struct {
	locationRepository repository.LocationRepository
}

type LocationService interface {
	GetLocation(ctx context.Context, uuid model.UUID) (*model.Location, error)
}

func NewLocationService(locationRepository repository.LocationRepository) LocationService {
	return &locationService{
		locationRepository,
	}
}

func (s *locationService) GetLocation(ctx context.Context, uuid model.UUID) (*model.Location, error) {
	return s.locationRepository.Find(ctx, string(uuid))
}
