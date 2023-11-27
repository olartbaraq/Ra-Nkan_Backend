-- name: CreateCategory :one
INSERT INTO category (
    name
) VALUES (
    $1) RETURNING *;

-- name: GetCategoryById :one
SELECT * FROM category WHERE id = $1;

-- name: GetCategoryByName :one
SELECT * FROM category WHERE name = $1;

-- name: ListAllCategory :many
SELECT * FROM category ORDER BY id LIMIT $1 OFFSET $2;

-- name: UpdateCategory :one
UPDATE category SET name = $2, updated_at = $3 WHERE id = $1 RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM category WHERE id = $1;