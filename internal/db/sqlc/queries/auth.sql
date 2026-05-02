-- name: CreateUser :one
INSERT INTO users (
    email, 
    password_hash
) VALUES (
    $1, $2
)
RETURNING id;

-- name: GetUserByEmail :one
SELECT id, email, password_hash
FROM users
WHERE email = $1 LIMIT 1;

-- name: ExistsUserByID :one
SELECT EXISTS (
    SELECT 1 FROM users WHERE id = $1
);