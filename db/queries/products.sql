-- name: CreateProduct :one
INSERT INTO products (
    name,
    description,
    price,
    image,
    qty_aval,
    shop_id
) VALUES (
    $1, $2, $3, $4, $5, $6) RETURNING *;

-- name: GetProductById :one
SELECT * FROM products WHERE id = $1;

-- name: GetProductByName :many
SELECT * FROM products WHERE name = $1 ORDER BY id;

-- name: GetProductByShop :many
SELECT * FROM products WHERE shop_id = $1 ORDER BY id;

-- name: ListAllProduct :many
SELECT * FROM products ORDER BY id LIMIT $1 OFFSET $2;

-- name: UpdateProduct :one
UPDATE products SET name = $2, qty_aval = $6, description = $5, price = $4, image = $3, updated_at = $6 WHERE id = $1 RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1;