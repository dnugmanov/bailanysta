-- name: CreatePost :one
INSERT INTO posts (author_id, text, course_id, module_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetPostByID :one
SELECT * FROM posts WHERE id = $1;

-- name: GetPostWithCounts :one
SELECT p.*,
       COUNT(DISTINCT l.user_id) as like_count,
       COUNT(DISTINCT c.id) as comment_count
FROM posts p
LEFT JOIN likes l ON p.id = l.post_id
LEFT JOIN comments c ON p.id = c.post_id
WHERE p.id = $1
GROUP BY p.id;

-- name: UpdatePost :one
UPDATE posts
SET text = $2, course_id = $3, module_id = $4, updated_at = now()
WHERE id = $1 AND author_id = $2
RETURNING *;

-- name: DeletePost :exec
DELETE FROM posts WHERE id = $1 AND author_id = $2;

-- name: GetPostsByAuthor :many
SELECT * FROM posts
WHERE author_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetFeedPosts :many
SELECT p.*,
       COUNT(DISTINCT l.user_id) as like_count,
       COUNT(DISTINCT c.id) as comment_count,
       u.username, u.avatar_url
FROM posts p
JOIN users u ON p.author_id = u.id
LEFT JOIN likes l ON p.id = l.post_id
LEFT JOIN comments c ON p.id = c.post_id
WHERE p.author_id IN (
    SELECT followee_id FROM follows WHERE follower_id = $1
) OR p.author_id = $1
GROUP BY p.id, u.username, u.avatar_url
ORDER BY p.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetPopularPosts :many
SELECT p.*,
       COUNT(DISTINCT l.user_id) as like_count,
       COUNT(DISTINCT c.id) as comment_count,
       u.username, u.avatar_url
FROM posts p
JOIN users u ON p.author_id = u.id
LEFT JOIN likes l ON p.id = l.post_id
LEFT JOIN comments c ON p.id = c.post_id
GROUP BY p.id, u.username, u.avatar_url
ORDER BY like_count DESC, p.created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetPostsByCourse :many
SELECT p.*,
       COUNT(DISTINCT l.user_id) as like_count,
       COUNT(DISTINCT c.id) as comment_count,
       u.username, u.avatar_url
FROM posts p
JOIN users u ON p.author_id = u.id
LEFT JOIN likes l ON p.id = l.post_id
LEFT JOIN comments c ON p.id = c.post_id
WHERE p.course_id = $1
GROUP BY p.id, u.username, u.avatar_url
ORDER BY p.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetPostsByModule :many
SELECT p.*,
       COUNT(DISTINCT l.user_id) as like_count,
       COUNT(DISTINCT c.id) as comment_count,
       u.username, u.avatar_url
FROM posts p
JOIN users u ON p.author_id = u.id
LEFT JOIN likes l ON p.id = l.post_id
LEFT JOIN comments c ON p.id = c.post_id
WHERE p.module_id = $1
GROUP BY p.id, u.username, u.avatar_url
ORDER BY p.created_at DESC
LIMIT $2 OFFSET $3;

-- name: SearchPosts :many
SELECT p.*,
       COUNT(DISTINCT l.user_id) as like_count,
       COUNT(DISTINCT c.id) as comment_count,
       u.username, u.avatar_url
FROM posts p
JOIN users u ON p.author_id = u.id
LEFT JOIN likes l ON p.id = l.post_id
LEFT JOIN comments c ON p.id = c.post_id
WHERE p.text ILIKE '%' || $1 || '%'
GROUP BY p.id, u.username, u.avatar_url
ORDER BY p.created_at DESC
LIMIT $2 OFFSET $3;
