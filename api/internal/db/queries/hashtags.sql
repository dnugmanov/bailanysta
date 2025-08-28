-- name: CreateHashtag :one
INSERT INTO hashtags (tag)
VALUES ($1)
ON CONFLICT (tag) DO UPDATE SET tag = EXCLUDED.tag
RETURNING *;

-- name: GetHashtagByTag :one
SELECT * FROM hashtags WHERE tag = $1;

-- name: GetHashtagByID :one
SELECT * FROM hashtags WHERE id = $1;

-- name: GetAllHashtags :many
SELECT * FROM hashtags ORDER BY tag;

-- name: CreatePostHashtag :exec
INSERT INTO post_hashtags (post_id, hashtag_id)
VALUES ($1, $2)
ON CONFLICT (post_id, hashtag_id) DO NOTHING;

-- name: DeletePostHashtags :exec
DELETE FROM post_hashtags WHERE post_id = $1;

-- name: GetHashtagsByPost :many
SELECT h.*
FROM hashtags h
JOIN post_hashtags ph ON h.id = ph.hashtag_id
WHERE ph.post_id = $1
ORDER BY h.tag;

-- name: GetPostsByHashtag :many
SELECT p.*,
       COUNT(DISTINCT l.user_id) as like_count,
       COUNT(DISTINCT c.id) as comment_count,
       u.username, u.avatar_url
FROM posts p
JOIN users u ON p.author_id = u.id
JOIN post_hashtags ph ON p.id = ph.post_id
JOIN hashtags h ON ph.hashtag_id = h.id
LEFT JOIN likes l ON p.id = l.post_id
LEFT JOIN comments c ON p.id = c.post_id
WHERE h.tag = $1
GROUP BY p.id, u.username, u.avatar_url
ORDER BY p.created_at DESC
LIMIT $2 OFFSET $3;
