package service

import (
	"context"

	"github.com/YuraSahanovskyi/booking-system/internal/domain"
	"github.com/google/uuid"
)

type AuthService interface {
	Register(ctx context.Context, email, password string) (uuid.UUID, error)
	Login(ctx context.Context, email, password string) (string, error) // повертає JWT токен
	ParseToken(ctx context.Context, token string) (uuid.UUID, error)
}

type EventService interface {
	GetAllEvents(ctx context.Context) ([]domain.Event, error)
	GetEventWithSeats(ctx context.Context, eventID uuid.UUID) (*domain.Event, []domain.Seat, error)
}

type BookingService interface {
	BookSeat(ctx context.Context, userID, seatID uuid.UUID) (*domain.Booking, error)
	GetUserBookings(ctx context.Context, userID uuid.UUID) ([]domain.Booking, error)
	ConfirmPayment(ctx context.Context, orderID uuid.UUID) error
	CancelBooking(ctx context.Context, userID, bookingID uuid.UUID) error
	CleanupExpiredBookings(ctx context.Context) (int64, error)
}
