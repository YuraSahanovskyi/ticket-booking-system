package worker_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/YuraSahanovskyi/booking-system/internal/domain"
	"github.com/YuraSahanovskyi/booking-system/internal/worker"
)


type mockBookingService struct {
	cleanupCalledCount int32
}

func (m *mockBookingService) CleanupExpiredBookings(ctx context.Context) (int64, error) {
	atomic.AddInt32(&m.cleanupCalledCount, 1)
	return 5, nil
}

func (m *mockBookingService) BookSeat(ctx context.Context, userID, seatID uuid.UUID) (*domain.Booking, error) {
	return nil, nil
}
func (m *mockBookingService) GetUserBookings(ctx context.Context, userID uuid.UUID) ([]domain.BookingWithDetails, error) {
	return nil, nil
}
func (m *mockBookingService) CancelBooking(ctx context.Context, userID, bookingID uuid.UUID) error {
	return nil
}
func (m *mockBookingService) ConfirmPayment(ctx context.Context, orderID uuid.UUID) error {
	return nil
}


func TestCleanupWorker_Lifecycle(t *testing.T) {
	mockSvc := &mockBookingService{}

	interval := 2 * time.Millisecond
	w := worker.NewCleanupWorker(mockSvc, interval)

	ctx, cancel := context.WithCancel(context.Background())

	go w.Start(ctx)

	time.Sleep(10 * time.Millisecond)

	cancel()

	time.Sleep(2 * time.Millisecond)

	calls := atomic.LoadInt32(&mockSvc.cleanupCalledCount)
	assert.GreaterOrEqual(t, calls, int32(1), "Воркер мав виконати хоча б одне очищення")
}
