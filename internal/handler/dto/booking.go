package dto

import (
	"time"

	"github.com/YuraSahanovskyi/booking-system/internal/domain"
)

// request POST /bookings
type CreateBookingRequest struct {
	SeatID string `json:"seat_id" binding:"required,uuid4" example:"9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d"`
}

// response POST /bookings
type CreateBookingResponse struct {
	ID        string    `json:"id" example:"a1b2c3d4..."`
	Status    string    `json:"status" example:"reserved"`
	ExpiresAt time.Time `json:"expires_at" example:"2026-05-02T18:15:00Z"`
}

func ToCreateBookingResponse(booking domain.Booking) CreateBookingResponse {
	return CreateBookingResponse{
		ID:        booking.ID.String(),
		Status:    booking.Status,
		ExpiresAt: booking.ExpiresAt,
	}
}

// response GET /bookings
type BookingResponse struct {
	ID             string    `json:"id" example:"a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d"`
	Status         string    `json:"status" example:"reserved"`
	ExpiresAt      time.Time `json:"expires_at" example:"2026-05-02T18:00:00Z"`
	EventTitle     string    `json:"event_title" example:"Kyiv Tech Summit"`
	EventStartTime time.Time `json:"event_start_time" example:"2026-05-20T18:00:00Z"`
	SeatRow        int       `json:"seat_row" example:"5"`
	SeatNumber     int       `json:"seat_number" example:"12"`
}

func ToBookingsResponse(bookings []domain.BookingWithDetails) []BookingResponse {
	res := make([]BookingResponse, len(bookings))
	for i, b := range bookings {
		res[i] = BookingResponse{
			ID:             b.ID.String(),
			Status:         b.Status,
			ExpiresAt:      b.ExpiresAt,
			EventTitle:     b.EventTitle,
			EventStartTime: b.EventStartTime,
			SeatRow:        b.SeatRow,
			SeatNumber:     b.SeatNumber,
		}
	}
	return res
}
