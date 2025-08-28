-- name: CreateLike :exec
INSERT INTO likes (user_id, post_id)
VALUES ($1, $2)
ON CONFLICT (user_id, post_id) DO NOTHING;

-- name: DeleteLike :exec
DELETE FROM likes WHERE user_id = $1 AND post_id = $2;

-- name: GetLikeByUserAndPost :one
SELECT * FROM likes WHERE user_id = $1 AND post_id = $2;

-- name: GetLikesByPost :many
SELECT l.*, u.username, u.avatar_url
FROM likes l
JOIN users u ON l.user_id = u.id
WHERE l.post_id = $1
ORDER BY l.created_at DESC;

-- name: GetLikesCount :one
SELECT COUNT(*) as count FROM likes WHERE post_id = $1;

-- name: GetLikesByUser :many
SELECT l.*, p.text as post_text, p.created_at as post_created_at
FROM likes l
JOIN posts p ON l.post_id = p.id
WHERE l.user_id = $1
ORDER BY l.created_at DESC
LIMIT $2 OFFSET $3;
