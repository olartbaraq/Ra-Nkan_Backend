-- name: CreateSubCategory :one
INSERT INTO sub_category (
    name,
    category_id
) VALUES (
    $1, $2) RETURNING *;

-- name: GetSubCategoryById :one
SELECT * FROM sub_category WHERE id = $1;

-- name: GetSubCategoryByName :one
SELECT * FROM sub_category WHERE name = $1;

-- name: GetSubCategoryByCategory :many
SELECT * FROM sub_category WHERE category_id = $1 ORDER BY id;

-- name: ListAllSubCategory :many
SELECT * FROM sub_category ORDER BY id LIMIT $1 OFFSET $2;

-- name: UpdateSubCategory :one
UPDATE sub_category SET name = $2, updated_at = $3 WHERE id = $1 RETURNING *;

-- name: DeleteSubCategory :exec
DELETE FROM sub_category WHERE id = $1;