-- name: GetEvents :many
SELECT * FROM events
ORDER BY start_time ASC;

-- name: GetEventByID :one
SELECT * FROM events
WHERE id = $1 LIMIT 1;

-- name: GetSeatsByEventWithBookings :many
SELECT 
    s.id AS seat_id, 
    s.event_id, 
    s.row, 
    s.number, 
    s.price,
    b.id AS booking_id, 
    b.user_id AS booking_user_id, 
    b.status AS booking_status, 
    b.expires_at AS booking_expires_at
FROM seats s
LEFT JOIN bookings b ON s.id = b.seat_id 
    AND b.status != 'canceled'
WHERE s.event_id = $1
ORDER BY s.row, s.number;