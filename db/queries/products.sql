-- name: CreateProduct :one
INSERT INTO products (
    name,
    description,
    price,
    images,
    qty_aval,
    shop_id,
    shop_name,
    category_id,
    category_name,
    sub_category_id,
    sub_category_name
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING *;

-- name: GetProductById :one
SELECT * FROM products WHERE id = $1;

-- name: GetProductByName :many
SELECT * FROM products WHERE name LIKE '%' || $1 || '%' ORDER BY id;

-- name: GetProductByShop :many
SELECT * FROM products WHERE shop_name = $1 ORDER BY id;

-- name: GetProductByPrice :many
SELECT * FROM products WHERE price = $1 ORDER BY id;

-- name: GetProductByPCS :many
SELECT * FROM products WHERE price = $1 AND sub_category_id = $2 AND category_id = $3 ORDER BY id;

-- name: GetProductBySubCategory :many
SELECT * FROM products WHERE sub_category_id = $1 ORDER BY id;

-- name: GetProductByCategory :many
SELECT * FROM products WHERE category_id = $1 ORDER BY id;

-- name: ListAllProducts :many
SELECT * FROM products ORDER BY id LIMIT $1 OFFSET $2;

-- name: ListAllProductsByOrders :many
SELECT
    p.id AS product_id,
    p.name AS product_name,
    COUNT(o.id) AS order_count
FROM
    products p
LEFT JOIN
    orders o ON p.id = o.product_id
GROUP BY
    p.id, p.name
ORDER BY
    order_count DESC;


-- name: UpdateProduct :one
UPDATE products SET name = $2, qty_aval = $6, description = $5, price = $4, images = $3, updated_at = $7 WHERE id = $1 RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1;

-- name: DeleteAllProducts :exec
DELETE FROM products;