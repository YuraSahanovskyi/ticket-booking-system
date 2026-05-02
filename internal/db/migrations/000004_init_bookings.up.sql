CREATE TYPE booking_status AS ENUM (
    'reserved',
    'paid',
    'canceled'
);

CREATE TABLE bookings (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    seat_id UUID NOT NULL REFERENCES seats(id),

    status booking_status NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX unique_active_booking_per_seat
ON bookings(seat_id)
WHERE status IN ('reserved', 'paid');