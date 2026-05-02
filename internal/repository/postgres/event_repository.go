package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/YuraSahanovskyi/booking-system/internal/db/sqlc"
	"github.com/YuraSahanovskyi/booking-system/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type EventRepository struct {
	q *sqlc.Queries
}

func NewEventRepository(q *sqlc.Queries) *EventRepository {
	return &EventRepository{q: q}
}

func (r *EventRepository) GetList(ctx context.Context) ([]domain.Event, error) {
	events, err := r.q.GetEvents(ctx)
	if err != nil {
		return nil, err
	}

	domainEvents := make([]domain.Event, len(events))
	for i, e := range events {
		domainEvents[i] = domain.Event{
			ID:          e.ID,
			Title:       e.Title,
			Description: mapPointerString(e.Description),
			Location:    mapPointerString(e.Location),
			StartTime:   e.StartTime,
			EndTime:     e.EndTime,
			CreatedAt:   e.CreatedAt,
		}
	}
	return domainEvents, nil
}

func (r *EventRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	e, err := r.q.GetEventByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrEventNotFound
		}
		return nil, err
	}

	return &domain.Event{
		ID:          e.ID,
		Title:       e.Title,
		Description: mapPointerString(e.Description),
		Location:    mapPointerString(e.Location),
		StartTime:   e.StartTime,
		EndTime:     e.EndTime,
		CreatedAt:   e.CreatedAt,
	}, nil
}

func (r *EventRepository) GetSeats(ctx context.Context, eventID uuid.UUID) ([]domain.Seat, error) {
	rows, err := r.q.GetSeatsByEventWithBookings(ctx, eventID)
	if err != nil {
		return nil, err
	}

	domainSeats := make([]domain.Seat, len(rows))
	for i, row := range rows {
		domainSeats[i] = domain.Seat{
			ID:      row.SeatID,
			EventID: row.EventID,
			Row:     int(row.Row),
			Number:  int(row.Number),
			Price:   int(row.Price),
		}

		if row.BookingID != uuid.Nil {
			var expiresAt time.Time
			if row.BookingExpiresAt != nil {
				expiresAt = *row.BookingExpiresAt
			}

			domainSeats[i].Booking = &domain.Booking{
				ID:        row.BookingID,
				UserID:    row.BookingUserID,
				Status:    mapStatus(row.BookingStatus),
				ExpiresAt: expiresAt,
			}
		}
	}

	return domainSeats, nil
}

func mapStatus(status sqlc.NullBookingStatus) string {
	if status.Valid {
		return string(status.BookingStatus)
	}
	return ""
}

func mapPointerString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
