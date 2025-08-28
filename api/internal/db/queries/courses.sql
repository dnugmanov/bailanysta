-- name: CreateCourse :one
INSERT INTO courses (title, description)
VALUES ($1, $2)
RETURNING *;

-- name: GetCourseByID :one
SELECT * FROM courses WHERE id = $1;

-- name: GetAllCourses :many
SELECT * FROM courses ORDER BY title;

-- name: UpdateCourse :one
UPDATE courses
SET title = $2, description = $3
WHERE id = $1
RETURNING *;

-- name: DeleteCourse :exec
DELETE FROM courses WHERE id = $1;

-- name: CreateModule :one
INSERT INTO modules (course_id, title, "order")
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetModuleByID :one
SELECT * FROM modules WHERE id = $1;

-- name: GetModulesByCourse :many
SELECT * FROM modules
WHERE course_id = $1
ORDER BY "order" ASC;

-- name: UpdateModule :one
UPDATE modules
SET title = $2, "order" = $3
WHERE id = $1
RETURNING *;

-- name: DeleteModule :exec
DELETE FROM modules WHERE id = $1;
