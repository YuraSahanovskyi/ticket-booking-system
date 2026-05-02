-- name: CreateBooking :one
INSERT INTO bookings (user_id, seat_id, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetBookingByID :one
SELECT * FROM bookings
WHERE id = $1 AND user_id = $2 LIMIT 1;

-- name: GetBookingsByUserID :many
SELECT * FROM bookings
WHERE user_id = $1
ORDER BY created_at DESC;

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