-- name: CreateUser :one
INSERT INTO users (
    lastname,
    firstname,
    email,
    phone,
    address,
    hashed_password,
    is_admin
) VALUES (
    $1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: ListAllUsers :many
SELECT * FROM users ORDER BY id LIMIT $1 OFFSET $2;

-- name: UpdateUserPassword :one
UPDATE users SET hashed_password = $1 WHERE id = $2 RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;