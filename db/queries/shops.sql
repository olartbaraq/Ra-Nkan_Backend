-- name: CreateShop :one
INSERT INTO shops (
    name,
    phone,
    address,
    email
) VALUES (
    $1, $2, $3, $4) RETURNING *;

-- name: GetShopByname :one
SELECT * FROM shops WHERE name = $1;

-- name: GetShopByEmail :one
SELECT * FROM shops WHERE email = $1;

-- name: ListAllShops :many
SELECT * FROM shops ORDER BY id LIMIT $1 OFFSET $2;

-- name: UpdateShop :one
UPDATE shops SET name = $2, address = $5, phone = $4, email = $3, updated_at = $6 WHERE id = $1 RETURNING *;

-- name: DeleteShop :exec
DELETE FROM shops WHERE id = $1;

-- name: DeleteAllShops :exec
DELETE FROM shops;