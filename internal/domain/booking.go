package domain

import (
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	SeatID    uuid.UUID
	Status    string
	ExpiresAt time.Time
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