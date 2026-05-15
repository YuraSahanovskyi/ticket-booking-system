-- name: CreateBooking :one
INSERT INTO bookings (user_id, seat_id, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetBookingByID :one
SELECT * FROM bookings
WHERE id = $1 AND user_id = $2 LIMIT 1;

-- name: GetBookingsByUserID :many
SELECT 
    b.id AS booking_id, 
    b.status, 
    b.expires_at,
    e.title AS event_title, 
    e.start_time AS event_start_time,
    s.row AS seat_row, 
    s.number AS seat_number
FROM bookings b
JOIN seats s ON b.seat_id = s.id
JOIN events e ON s.event_id = e.id
WHERE b.user_id = $1
ORDER BY e.start_time DESC;

-- name: SetBookingStatusPaid :exec
UPDATE bookings
SET status = 'paid'
WHERE id = $1 AND status = 'reserved';

-- name: CancelBooking :exec
UPDATE bookings
SET status = 'canceled'
WHERE id = $1 AND user_id = $2 AND status = 'reserved';

-- name: CancelExpiredBookings :execrows
UPDATE bookings
SET status = 'canceled'
WHERE status = 'reserved' 
  AND expires_at < NOW();