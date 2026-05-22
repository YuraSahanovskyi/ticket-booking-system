package service

import (
	"context"

	"github.com/YuraSahanovskyi/booking-system/internal/domain"
	"github.com/YuraSahanovskyi/booking-system/internal/repository"
	"github.com/google/uuid"
)

type eventService struct {
	repo repository.EventRepository
}

func NewEventService(repo repository.EventRepository) EventService {
	return &eventService{
		repo: repo,
	}
}

func (s *eventService) GetAllEvents(ctx context.Context) ([]domain.Event, error) {
	return s.repo.GetList(ctx)
}

func (s *eventService) GetEventWithSeats(ctx context.Context, eventID uuid.UUID) (*domain.Event, []domain.Seat, error) {
	event, err := s.repo.GetByID(ctx, eventID)
	if err != nil {
		return nil, nil, err
	}

	seats, err := s.repo.GetSeats(ctx, eventID)
	if err != nil {
		return nil, nil, err
	}

	return event, seats, nil
}
