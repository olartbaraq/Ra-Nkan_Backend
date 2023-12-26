-- name: CreateOrder :one
INSERT INTO orders (
    product_id,
    qty_bought,
    unit_price,
    total_price,
    user_id
) VALUES (
    $1, $2, $3, $4, $5) RETURNING *;