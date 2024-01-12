-- name: CreateOrder :one
INSERT INTO orders (
    user_id,
    items
) VALUES (
    $1, $2) RETURNING *;

-- name: GetOrderById :one
SELECT * FROM orders WHERE id = $1;


-- name: GetOrdersByUser :many
select * FROM orders WHERE user_id = $1;