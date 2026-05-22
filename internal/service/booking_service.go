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
}

func NewBookingService(br repository.BookingRepository, er repository.EventRepository, ttl time.Duration, rdb *redis.Client) BookingService {
	return &bookingService{
		bookingRepo: br,
		eventRepo:   er,
		bookingTTL:  ttl,
		rdb:         rdb,
	}
}

func (s *bookingService) BookSeat(ctx context.Context, userID, seatID uuid.UUID) (*domain.Booking, error) {
	statusKey := fmt.Sprintf("seat:status:%s", seatID.String())

	cachedStatus, err := s.rdb.Get(ctx, statusKey).Result()
	if err == nil && cachedStatus == "occupied" {
		return nil, errors.New("this seat is already booked or reserved (cached)")
	}

	lockKey := fmt.Sprintf("lock:seat:%s", seatID.String())
	lockValue := uuid.New().String()

	success, err := s.rdb.SetNX(ctx, lockKey, lockValue, 5*time.Second).Result()
	if err != nil {
		return nil, fmt.Errorf("redis lock error: %w", err)
	}
	if !success {
		return nil, errors.New("seat is temporarily locked, please try again")
	}

	defer func() {
		var luaReleaseLock = redis.NewScript(`
			if redis.call("get", KEYS[1]) == ARGV[1] then
				return redis.call("del", KEYS[1])
			else
				return 0
			end
		`)
		_ = luaReleaseLock.Run(ctx, s.rdb, []string{lockKey}, lockValue).Err()
	}()

	booking := &domain.Booking{
		UserID:    userID,
		SeatID:    seatID,
		Status:    "reserved",
		ExpiresAt: time.Now().Add(s.bookingTTL),
	}

	createdBooking, err := s.bookingRepo.Create(ctx, booking)
	if err != nil {
		_ = s.rdb.Set(ctx, statusKey, "occupied", 2*time.Minute).Err()
		return nil, err
	}

	_ = s.rdb.Set(ctx, statusKey, "occupied", 2*time.Minute).Err()

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
