package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/YuraSahanovskyi/booking-system/internal/domain"
	"github.com/YuraSahanovskyi/booking-system/internal/handler"
)

type mockEventService struct {
	getAllEventsFunc       func(ctx context.Context) ([]domain.Event, error)
	getEventWithSeatsFunc func(ctx context.Context, eventID uuid.UUID) (*domain.Event, []domain.Seat, error)
}

func (m *mockEventService) GetAllEvents(ctx context.Context) ([]domain.Event, error) {
	if m.getAllEventsFunc != nil {
		return m.getAllEventsFunc(ctx)
	}
	return nil, nil
}

func (m *mockEventService) GetEventWithSeats(ctx context.Context, eventID uuid.UUID) (*domain.Event, []domain.Seat, error) {
	if m.getEventWithSeatsFunc != nil {
		return m.getEventWithSeatsFunc(ctx, eventID)
	}
	return &domain.Event{}, nil, nil
}

func TestEventHandler_GetAllEvents(t *testing.T) {
	t.Run("200 OK - Get List of Events", func(t *testing.T) {
		mockEvent := &mockEventService{
			getAllEventsFunc: func(ctx context.Context) ([]domain.Event, error) {
				return []domain.Event{
					{ID: uuid.New(), Title: "Concert 1"},
					{ID: uuid.New(), Title: "Concert 2"},
				}, nil
			},
		}

		h := handler.NewHandler(nil, mockEvent, nil)
		router := h.Init()

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/events/", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "Concert 1")
		assert.Contains(t, rr.Body.String(), "Concert 2")
	})
}

func TestEventHandler_GetEventSeats(t *testing.T) {
	t.Run("200 OK - Get Event Seats", func(t *testing.T) {
		eventID := uuid.New()
		mockEvent := &mockEventService{
			getEventWithSeatsFunc: func(ctx context.Context, id uuid.UUID) (*domain.Event, []domain.Seat, error) {
				assert.Equal(t, eventID, id)
				return &domain.Event{ID: eventID, Title: "Epic Show"}, []domain.Seat{
					{ID: uuid.New(), EventID: eventID, Row: 5, Number: 12},
				}, nil
			},
		}

		h := handler.NewHandler(nil, mockEvent, nil)
		router := h.Init()

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/events/"+eventID.String()+"/seats", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "Epic Show")
	})

	t.Run("400 Bad Request - Invalid Event ID UUID", func(t *testing.T) {
		h := handler.NewHandler(nil, nil, nil)
		router := h.Init()

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/events/not-a-valid-uuid/seats", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "invalid id")
	})
}