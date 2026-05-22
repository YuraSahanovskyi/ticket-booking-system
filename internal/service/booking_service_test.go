package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	"github.com/YuraSahanovskyi/booking-system/internal/domain"
	"github.com/YuraSahanovskyi/booking-system/internal/service"
)

type mockBookingRepo struct {
	createFunc                func(ctx context.Context, b *domain.Booking) (*domain.Booking, error)
	getByUserIDFunc           func(ctx context.Context, userID uuid.UUID) ([]domain.BookingWithDetails, error)
	getByIDFunc               func(ctx context.Context, bookingID, userID uuid.UUID) (*domain.Booking, error)
	cancelFunc                func(ctx context.Context, bookingID, userID uuid.UUID) error
	updateStatusToPaidFunc    func(ctx context.Context, orderID uuid.UUID) error
	cancelExpiredBookingsFunc func(ctx context.Context) (int64, error)
}

func (m *mockBookingRepo) Create(ctx context.Context, b *domain.Booking) (*domain.Booking, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, b)
	}
	return b, nil
}

func (m *mockBookingRepo) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.BookingWithDetails, error) {
	if m.getByUserIDFunc != nil {
		return m.getByUserIDFunc(ctx, userID)
	}
	return nil, nil
}

func (m *mockBookingRepo) GetByID(ctx context.Context, bookingID, userID uuid.UUID) (*domain.Booking, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, bookingID, userID)
	}
	return &domain.Booking{}, nil
}

func (m *mockBookingRepo) Cancel(ctx context.Context, bookingID, userID uuid.UUID) error {
	if m.cancelFunc != nil {
		return m.cancelFunc(ctx, bookingID, userID)
	}
	return nil
}

func (m *mockBookingRepo) UpdateStatusToPaid(ctx context.Context, orderID uuid.UUID) error {
	if m.updateStatusToPaidFunc != nil {
		return m.updateStatusToPaidFunc(ctx, orderID)
	}
	return nil
}

func (m *mockBookingRepo) CancelExpiredBookings(ctx context.Context) (int64, error) {
	if m.cancelExpiredBookingsFunc != nil {
		return m.cancelExpiredBookingsFunc(ctx)
	}
	return 0, nil
}


func TestBookingService(t *testing.T) {
	mr, err := miniredis.Run()
	assert.NoError(t, err)
	defer mr.Close()

	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})

	t.Run("BookSeat - Success", func(t *testing.T) {
		mockRepo := &mockBookingRepo{
			createFunc: func(ctx context.Context, b *domain.Booking) (*domain.Booking, error) {
				return b, nil
			},
		}
		svc := service.NewBookingService(mockRepo, nil, 15*time.Minute, rdb, true, true)

		res, err := svc.BookSeat(context.Background(), uuid.New(), uuid.New())
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "reserved", res.Status)
	})

	t.Run("GetUserBookings - Success", func(t *testing.T) {
		userID := uuid.New()
		expectedBookings := []domain.BookingWithDetails{
			{
				ID:             uuid.New(),
				Status:         "reserved",
				ExpiresAt:      time.Now().Add(15 * time.Minute),
				EventTitle:     "Test Concert",
				EventStartTime: time.Now().Add(24 * time.Hour),
				SeatRow:        1,
				SeatNumber:     5,
			},
		}

		mockRepo := &mockBookingRepo{
			getByUserIDFunc: func(ctx context.Context, id uuid.UUID) ([]domain.BookingWithDetails, error) {
				assert.Equal(t, userID, id)
				return expectedBookings, nil
			},
		}
		svc := service.NewBookingService(mockRepo, nil, 15*time.Minute, rdb, true, true)

		res, err := svc.GetUserBookings(context.Background(), userID)
		assert.NoError(t, err)
		assert.Len(t, res, 1)
		assert.Equal(t, expectedBookings[0].ID, res[0].ID)
		assert.Equal(t, "Test Concert", res[0].EventTitle)
	})

	t.Run("CancelBooking - Success", func(t *testing.T) {
		userID := uuid.New()
		bookingID := uuid.New()

		mockRepo := &mockBookingRepo{
			getByIDFunc: func(ctx context.Context, bID, uID uuid.UUID) (*domain.Booking, error) {
				return &domain.Booking{ID: bookingID, UserID: userID, Status: "reserved"}, nil
			},
			cancelFunc: func(ctx context.Context, bID, uID uuid.UUID) error {
				assert.Equal(t, bookingID, bID)
				assert.Equal(t, userID, uID)
				return nil
			},
		}
		svc := service.NewBookingService(mockRepo, nil, 15*time.Minute, rdb, true, true)

		err := svc.CancelBooking(context.Background(), userID, bookingID)
		assert.NoError(t, err)
	})

	t.Run("CancelBooking - Fail when not reserved", func(t *testing.T) {
		mockRepo := &mockBookingRepo{
			getByIDFunc: func(ctx context.Context, bID, uID uuid.UUID) (*domain.Booking, error) {
				return &domain.Booking{Status: "paid"}, nil
			},
		}
		svc := service.NewBookingService(mockRepo, nil, 15*time.Minute, rdb, true, true)

		err := svc.CancelBooking(context.Background(), uuid.New(), uuid.New())
		assert.ErrorIs(t, err, domain.ErrBookingCannotBeCanceled)
	})

	t.Run("ConfirmPayment - Success", func(t *testing.T) {
		orderID := uuid.New()
		called := false

		mockRepo := &mockBookingRepo{
			updateStatusToPaidFunc: func(ctx context.Context, id uuid.UUID) error {
				assert.Equal(t, orderID, id)
				called = true
				return nil
			},
		}
		svc := service.NewBookingService(mockRepo, nil, 15*time.Minute, rdb, true, true)

		err := svc.ConfirmPayment(context.Background(), orderID)
		assert.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("CleanupExpiredBookings - Success", func(t *testing.T) {
		var expectedCount int64 = 5

		mockRepo := &mockBookingRepo{
			cancelExpiredBookingsFunc: func(ctx context.Context) (int64, error) {
				return expectedCount, nil
			},
		}
		svc := service.NewBookingService(mockRepo, nil, 15*time.Minute, rdb, true, true)

		count, err := svc.CleanupExpiredBookings(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, count)
	})
}
