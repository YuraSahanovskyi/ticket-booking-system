package domain

import (
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	SeatID    uuid.UUID `json:"seat_id"`
	Status    string    `json:"status"`
	ExpiresAt time.Time `json:"expires_at"`
}

type BookingWithDetails struct {
	ID             uuid.UUID
	Status         string
	ExpiresAt      time.Time
	EventTitle     string
	EventStartTime time.Time
	SeatRow        int
	SeatNumber     int
}