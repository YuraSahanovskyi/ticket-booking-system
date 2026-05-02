package service

import (
	"context"
	"fmt"
	"time"

	"github.com/YuraSahanovskyi/booking-system/internal/domain"
	"github.com/YuraSahanovskyi/booking-system/internal/repository"
	"github.com/google/uuid"
)

type bookingService struct {
	bookingRepo repository.BookingRepository
	eventRepo   repository.EventRepository
	bookingTTL  time.Duration
}

func NewBookingService(br repository.BookingRepository, er repository.EventRepository, ttl time.Duration) BookingService {
	return &bookingService{
		bookingRepo: br,
		eventRepo:   er,
		bookingTTL:  ttl,
	}
}

func (s *bookingService) BookSeat(ctx context.Context, userID, seatID uuid.UUID) (*domain.Booking, error) {

	booking := &domain.Booking{
		UserID:    userID,
		SeatID:    seatID,
		Status:    "reserved",
		ExpiresAt: time.Now().Add(s.bookingTTL),
	}

	createdBooking, err := s.bookingRepo.Create(ctx, booking)
	if err != nil {
		return nil, err
	}

	return createdBooking, nil
}

func (s *bookingService) GetUserBookings(ctx context.Context, userID uuid.UUID) ([]domain.Booking, error) {
	return s.bookingRepo.GetByUserID(ctx, userID)
}

func (s *bookingService) CancelBooking(ctx context.Context, userID, bookingID uuid.UUID) error {
	b, err := s.bookingRepo.GetByID(ctx, bookingID, userID)
	if err != nil {
		return err
	}

	if b.Status != "reserved" {
		return domain.ErrBookingCannotBeCanceled
	}

	return s.bookingRepo.Cancel(ctx, bookingID, userID)
}

func (s *bookingService) ConfirmPayment(ctx context.Context, orderID uuid.UUID) error {
	err := s.bookingRepo.UpdateStatusToPaid(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to update booking status: %w", err)
	}

	return nil
}

func (s *bookingService) CleanupExpiredBookings(ctx context.Context) (int64, error) {
    return s.bookingRepo.CancelExpiredBookings(ctx)
}