package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/YuraSahanovskyi/booking-system/internal/domain"
	"github.com/YuraSahanovskyi/booking-system/internal/service"
)

type mockEventRepo struct {
	getListFunc  func(ctx context.Context) ([]domain.Event, error)
	getByIDFunc  func(ctx context.Context, id uuid.UUID) (*domain.Event, error)
	getSeatsFunc func(ctx context.Context, eventID uuid.UUID) ([]domain.Seat, error)
}

func (m *mockEventRepo) GetList(ctx context.Context) ([]domain.Event, error) {
	if m.getListFunc != nil {
		return m.getListFunc(ctx)
	}
	return nil, nil
}

func (m *mockEventRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return &domain.Event{}, nil
}

func (m *mockEventRepo) GetSeats(ctx context.Context, eventID uuid.UUID) ([]domain.Seat, error) {
	if m.getSeatsFunc != nil {
		return m.getSeatsFunc(ctx, eventID)
	}
	return nil, nil
}

func TestEventService(t *testing.T) {
	t.Run("GetAllEvents - Success", func(t *testing.T) {
		expectedEvents := []domain.Event{
			{ID: uuid.New(), Title: "Concert A"},
			{ID: uuid.New(), Title: "Concert B"},
		}

		mockRepo := &mockEventRepo{
			getListFunc: func(ctx context.Context) ([]domain.Event, error) {
				return expectedEvents, nil
			},
		}

		svc := service.NewEventService(mockRepo)
		res, err := svc.GetAllEvents(context.Background())

		assert.NoError(t, err)
		assert.Len(t, res, 2)
		assert.Equal(t, "Concert A", res[0].Title)
	})

	t.Run("GetEventWithSeats - Success", func(t *testing.T) {
		eventID := uuid.New()
		expectedEvent := &domain.Event{ID: eventID, Title: "Rock Show"}
		expectedSeats := []domain.Seat{
			{ID: uuid.New(), EventID: eventID, Row: 1, Number: 10},
		}

		mockRepo := &mockEventRepo{
			getByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
				assert.Equal(t, eventID, id)
				return expectedEvent, nil
			},
			getSeatsFunc: func(ctx context.Context, id uuid.UUID) ([]domain.Seat, error) {
				assert.Equal(t, eventID, id)
				return expectedSeats, nil
			},
		}

		svc := service.NewEventService(mockRepo)
		event, seats, err := svc.GetEventWithSeats(context.Background(), eventID)

		assert.NoError(t, err)
		assert.Equal(t, expectedEvent.Title, event.Title)
		assert.Len(t, seats, 1)
		assert.Equal(t, 10, seats[0].Number)
	})

	t.Run("GetEventWithSeats - Event Not Found", func(t *testing.T) {
		eventIDErr := errors.New("event not found")
		mockRepo := &mockEventRepo{
			getByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
				return nil, eventIDErr
			},
			getSeatsFunc: func(ctx context.Context, id uuid.UUID) ([]domain.Seat, error) {
				t.Fatal("GetSeats не мав викликатися, якщо івент не знайдено")
				return nil, nil
			},
		}

		svc := service.NewEventService(mockRepo)
		event, seats, err := svc.GetEventWithSeats(context.Background(), uuid.New())

		assert.ErrorIs(t, err, eventIDErr)
		assert.Nil(t, event)
		assert.Nil(t, seats)
	})

	t.Run("GetEventWithSeats - GetSeats Error", func(t *testing.T) {
		seatsErr := errors.New("database connection failure")
		mockRepo := &mockEventRepo{
			getByIDFunc: func(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
				return &domain.Event{Title: "Comedy Club"}, nil
			},
			getSeatsFunc: func(ctx context.Context, id uuid.UUID) ([]domain.Seat, error) {
				return nil, seatsErr
			},
		}

		svc := service.NewEventService(mockRepo)
		event, seats, err := svc.GetEventWithSeats(context.Background(), uuid.New())

		assert.ErrorIs(t, err, seatsErr)
		assert.Nil(t, event)
		assert.Nil(t, seats)
	})
}
