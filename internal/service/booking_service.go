package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/YuraSahanovskyi/booking-system/internal/domain"
	"github.com/YuraSahanovskyi/booking-system/internal/repository"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type bookingService struct {
	bookingRepo repository.BookingRepository
	eventRepo   repository.EventRepository
	bookingTTL  time.Duration
	rdb         *redis.Client
	enableCache bool
}

func NewBookingService(br repository.BookingRepository, er repository.EventRepository, ttl time.Duration, rdb *redis.Client, enableCache bool) BookingService {
	return &bookingService{
		bookingRepo: br,
		eventRepo:   er,
		bookingTTL:  ttl,
		rdb:         rdb,
		enableCache: enableCache,
	}
}

func (s *bookingService) BookSeat(ctx context.Context, userID, seatID uuid.UUID) (*domain.Booking, error) {
	statusKey := fmt.Sprintf("seat:status:%s", seatID.String())

	if s.enableCache {
		acquired, err := s.rdb.SetNX(ctx, statusKey, "occupied", 5*time.Minute).Result()
		if err != nil {
			return nil, fmt.Errorf("redis cache error: %w", err)
		}
		if !acquired {
			return nil, domain.ErrSeatAlreadyBooked
		}
	}

	booking := &domain.Booking{
		UserID:    userID,
		SeatID:    seatID,
		Status:    "reserved",
		ExpiresAt: time.Now().Add(s.bookingTTL),
	}

	createdBooking, err := s.bookingRepo.Create(ctx, booking)
	if err != nil {
		if !errors.Is(err, domain.ErrSeatAlreadyBooked) {
			if s.enableCache {
				_ = s.rdb.Del(ctx, statusKey).Err()
			}
		}
		return nil, err
	}

	return createdBooking, nil
}

func (s *bookingService) GetUserBookings(ctx context.Context, userID uuid.UUID) ([]domain.BookingWithDetails, error) {
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
