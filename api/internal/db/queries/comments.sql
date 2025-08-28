-- name: CreateComment :one
INSERT INTO comments (post_id, author_id, text)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetCommentsByPostID :many
SELECT c.*, u.username, u.avatar_url
FROM comments c
JOIN users u ON c.author_id = u.id
WHERE c.post_id = $1
ORDER BY c.created_at ASC
LIMIT $2 OFFSET $3;

-- name: GetCommentByID :one
SELECT c.*, u.username, u.avatar_url
FROM comments c
JOIN users u ON c.author_id = u.id
WHERE c.id = $1;

-- name: UpdateComment :one
UPDATE comments
SET text = $2
WHERE id = $1 AND author_id = $3
RETURNING *;

-- name: DeleteComment :exec
DELETE FROM comments WHERE id = $1 AND author_id = $2;

-- name: GetCommentsCount :one
SELECT COUNT(*) as count FROM comments WHERE post_id = $1;
