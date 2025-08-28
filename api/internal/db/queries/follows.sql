-- name: CreateFollow :exec
INSERT INTO follows (follower_id, followee_id)
VALUES ($1, $2)
ON CONFLICT (follower_id, followee_id) DO NOTHING;

-- name: DeleteFollow :exec
DELETE FROM follows WHERE follower_id = $1 AND followee_id = $2;

-- name: GetFollowByUsers :one
SELECT * FROM follows WHERE follower_id = $1 AND followee_id = $2;

-- name: GetFollowers :many
SELECT f.*, u.username, u.avatar_url, u.bio
FROM follows f
JOIN users u ON f.follower_id = u.id
WHERE f.followee_id = $1
ORDER BY f.created_at DESC;

-- name: GetFollowing :many
SELECT f.*, u.username, u.avatar_url, u.bio
FROM follows f
JOIN users u ON f.followee_id = u.id
WHERE f.follower_id = $1
ORDER BY f.created_at DESC;

-- name: GetFollowersCount :one
SELECT COUNT(*) as count FROM follows WHERE followee_id = $1;

-- name: GetFollowingCount :one
SELECT COUNT(*) as count FROM follows WHERE follower_id = $1;
