-- name: CreateNotification :one
INSERT INTO notifications (user_id, type, entity_id, payload_json)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetNotificationsByUser :many
SELECT * FROM notifications
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetUnreadNotifications :many
SELECT * FROM notifications
WHERE user_id = $1 AND read_at IS NULL
ORDER BY created_at DESC;

-- name: MarkNotificationAsRead :exec
UPDATE notifications
SET read_at = now()
WHERE id = $1 AND user_id = $2;

-- name: MarkAllNotificationsAsRead :exec
UPDATE notifications
SET read_at = now()
WHERE user_id = $1 AND read_at IS NULL;

-- name: DeleteNotification :exec
DELETE FROM notifications WHERE id = $1 AND user_id = $2;

-- name: GetNotificationsCount :one
SELECT COUNT(*) as count FROM notifications WHERE user_id = $1 AND read_at IS NULL;
