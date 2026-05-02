package repository

import (
	"context"

	"github.com/YuraSahanovskyi/booking-system/internal/domain"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (uuid.UUID, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
}

type EventRepository interface {
	GetList(ctx context.Context) ([]domain.Event, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error)
	GetSeats(ctx context.Context, eventID uuid.UUID) ([]domain.Seat, error)
}

type BookingRepository interface {
	Create(ctx context.Context, b *domain.Booking) (*domain.Booking, error)
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Booking, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Booking, error)
	UpdateStatusToPaid(ctx context.Context, id uuid.UUID) error
	Cancel(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	CancelExpiredBookings(ctx context.Context) (int64, error)
}
