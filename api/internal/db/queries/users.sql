-- name: CreateUser :one
INSERT INTO users (username, email, password_hash, bio, avatar_url)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: UpdateUser :one
UPDATE users
SET bio = $2, avatar_url = $3, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: GetUsersByIDs :many
SELECT * FROM users WHERE id = ANY($1::uuid[]);

-- name: SearchUsers :many
SELECT * FROM users
WHERE username ILIKE '%' || $1 || '%' OR bio ILIKE '%' || $1 || '%'
ORDER BY username
LIMIT $2 OFFSET $3;
