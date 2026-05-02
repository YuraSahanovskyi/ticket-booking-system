package worker

import (
    "context"
    "log"
    "time"
    "github.com/YuraSahanovskyi/booking-system/internal/service"
)

type CleanupWorker struct {
    bookingService service.BookingService
    interval       time.Duration
}

func NewCleanupWorker(bs service.BookingService, interval time.Duration) *CleanupWorker {
    return &CleanupWorker{
        bookingService: bs,
        interval:       interval,
    }
}

func (w *CleanupWorker) Start(ctx context.Context) {
    ticker := time.NewTicker(w.interval)
    defer ticker.Stop()

    log.Printf("Cleanup worker started with interval %v", w.interval)

    for {
        select {
        case <-ticker.C:
            count, err := w.bookingService.CleanupExpiredBookings(ctx)
            if err != nil {
                log.Printf("Cleanup worker error: %v", err)
                continue
            }
            if count > 0 {
                log.Printf("Cleanup worker: canceled %d expired bookings", count)
            }
        case <-ctx.Done():
            log.Println("Cleanup worker stopping...")
            return
        }
    }
}