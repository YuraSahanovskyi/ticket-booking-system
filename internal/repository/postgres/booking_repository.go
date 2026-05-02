package postgres

import (
	"context"
	"errors"

	"github.com/YuraSahanovskyi/booking-system/internal/db/sqlc"
	"github.com/YuraSahanovskyi/booking-system/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type BookingRepository struct {
	q *sqlc.Queries
}

func NewBookingRepository(q *sqlc.Queries) *BookingRepository {
	return &BookingRepository{q: q}
}

func (r *BookingRepository) Create(ctx context.Context, b *domain.Booking) (*domain.Booking, error) {
	row, err := r.q.CreateBooking(ctx, sqlc.CreateBookingParams{
		UserID:    b.UserID,
		SeatID:    b.SeatID,
		ExpiresAt: b.ExpiresAt,
	})

	if err != nil {
		var pgErr *pgconn.PgError
		// unique_violation (23505)
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, domain.ErrSeatAlreadyBooked
		}
		return nil, err
	}

	return &domain.Booking{
		ID:        row.ID,
		UserID:    row.UserID,
		SeatID:    row.SeatID,
		Status:    string(row.Status),
		ExpiresAt: row.ExpiresAt,
	}, nil
}

func (r *BookingRepository) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Booking, error) {
	b, err := r.q.GetBookingByID(ctx, sqlc.GetBookingByIDParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrBookingNotFound
		}
		return nil, err
	}

	return &domain.Booking{
		ID:        b.ID,
		UserID:    b.UserID,
		SeatID:    b.SeatID,
		Status:    string(b.Status),
		ExpiresAt: b.ExpiresAt,
	}, nil
}

func (r *BookingRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Booking, error) {
	rows, err := r.q.GetBookingsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	bookings := make([]domain.Booking, len(rows))
	for i, b := range rows {
		bookings[i] = domain.Booking{
			ID:        b.ID,
			UserID:    b.UserID,
			SeatID:    b.SeatID,
			Status:    string(b.Status),
			ExpiresAt: b.ExpiresAt,
		}
	}
	return bookings, nil
}

func (r *BookingRepository) UpdateStatusToPaid(ctx context.Context, id uuid.UUID) error {
	return r.q.SetBookingStatusPaid(ctx, id)
}

func (r *BookingRepository) Cancel(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	return r.q.CancelBooking(ctx, sqlc.CancelBookingParams{
		ID:     id,
		UserID: userID,
	})
}

func (r *BookingRepository) CancelExpiredBookings(ctx context.Context) (int64, error) {
    return r.q.CancelExpiredBookings(ctx)
}