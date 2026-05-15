package dto

import (
	"time"

	"github.com/YuraSahanovskyi/booking-system/internal/domain"
)

// response GET /events
type EventResponse struct {
	ID          string     `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Title       string     `json:"title" example:"Title"`
	Description string     `json:"description" example:"Description"`
	Location    string     `json:"location" example:"Location"`
	StartTime   time.Time  `json:"start_time" example:"2026-05-20T18:00:00Z"`
	EndTime     *time.Time `json:"end_time,omitempty" example:"2026-05-20T22:00:00Z"`
}

// part of response GET /events/{id}/seats
type SeatResponse struct {
	ID        string `json:"id" example:"9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d"`
	Row       int    `json:"row" example:"1"`
	Number    int    `json:"number" example:"12"`
	Price     int    `json:"price" example:"500"`
	Available bool   `json:"is_available" example:"true"`
}

// response GET /events/{id}/seats
type EventWithSeatsResponse struct {
	ID          string         `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Title       string         `json:"title" example:"Title"`
	Description string         `json:"description" example:"Description"`
	Location    string         `json:"location" example:"Location"`
	StartTime   time.Time      `json:"start_time" example:"2026-05-20T18:00:00Z"`
	EndTime     *time.Time     `json:"end_time,omitempty" example:"2026-05-20T22:00:00Z"`
	Seats       []SeatResponse `json:"seats"`
}

func ToEventsResponse(events []domain.Event) []EventResponse {
	dtoEvents := make([]EventResponse, len(events))
	for i, event := range events {
		dtoEvents[i] = EventResponse{
			ID:          event.ID.String(),
			Title:       event.Title,
			Description: event.Description,
			Location:    event.Location,
			StartTime:   event.StartTime,
			EndTime:     event.EndTime,
		}
	}
	return dtoEvents
}

func ToEventWithSeatsResponse(event domain.Event, seats []domain.Seat) EventWithSeatsResponse {
	dtoSeats := make([]SeatResponse, len(seats))

	for i, seat := range seats {
		dtoSeats[i] = toSeatResponse(seat)
	}
	return EventWithSeatsResponse{
		ID:          event.ID.String(),
		Title:       event.Title,
		Description: event.Description,
		Location:    event.Location,
		StartTime:   event.StartTime,
		EndTime:     event.EndTime,
		Seats:       dtoSeats,
	}
}

func toSeatResponse(s domain.Seat) SeatResponse {
	return SeatResponse{
		ID:        s.ID.String(),
		Row:       s.Row,
		Number:    s.Number,
		Price:     s.Price,
		Available: s.Booking == nil,
	}
}
