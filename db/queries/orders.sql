-- name: CreateOrder :one
INSERT INTO orders (
    user_id,
    items,
    transaction_ref,
    total_price,
    pay_ref,
    status
) VALUES (
    $1, $2, $3, $4, $5, $6) RETURNING *;

-- name: GetOrderById :one
SELECT * FROM orders WHERE id = $1;

-- name: GetOrdersByUser :many
select * FROM orders WHERE user_id = $1;

-- name: CompleteOrder :one
UPDATE orders SET pay_ref = $2, status = $3 WHERE id = $1 RETURNING *;