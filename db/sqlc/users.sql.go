// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: users.sql

package db

import (
	"context"
	"time"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    lastname,
    firstname,
    email,
    phone,
    address,
    hashed_password
) VALUES (
    $1, $2, $3, $4, $5, $6) RETURNING id, lastname, firstname, hashed_password, phone, address, email, is_admin, created_at, updated_at
`

type CreateUserParams struct {
	Lastname       string `json:"lastname"`
	Firstname      string `json:"firstname"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	Address        string `json:"address"`
	HashedPassword string `json:"hashed_password"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Lastname,
		arg.Firstname,
		arg.Email,
		arg.Phone,
		arg.Address,
		arg.HashedPassword,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Lastname,
		&i.Firstname,
		&i.HashedPassword,
		&i.Phone,
		&i.Address,
		&i.Email,
		&i.IsAdmin,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteAllUsers = `-- name: DeleteAllUsers :exec
DELETE FROM users
`

func (q *Queries) DeleteAllUsers(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteAllUsers)
	return err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, lastname, firstname, hashed_password, phone, address, email, is_admin, created_at, updated_at FROM users WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Lastname,
		&i.Firstname,
		&i.HashedPassword,
		&i.Phone,
		&i.Address,
		&i.Email,
		&i.IsAdmin,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserById = `-- name: GetUserById :one
SELECT id, lastname, firstname, hashed_password, phone, address, email, is_admin, created_at, updated_at FROM users WHERE id = $1
`

func (q *Queries) GetUserById(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserById, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Lastname,
		&i.Firstname,
		&i.HashedPassword,
		&i.Phone,
		&i.Address,
		&i.Email,
		&i.IsAdmin,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listAllUsers = `-- name: ListAllUsers :many
SELECT id, lastname, firstname, hashed_password, phone, address, email, is_admin, created_at, updated_at FROM users ORDER BY id LIMIT $1 OFFSET $2
`

type ListAllUsersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListAllUsers(ctx context.Context, arg ListAllUsersParams) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, listAllUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Lastname,
			&i.Firstname,
			&i.HashedPassword,
			&i.Phone,
			&i.Address,
			&i.Email,
			&i.IsAdmin,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUser = `-- name: UpdateUser :one
UPDATE users SET address = $4, phone = $3, email = $2, updated_at = $5 WHERE id = $1 RETURNING id, lastname, firstname, hashed_password, phone, address, email, is_admin, created_at, updated_at
`

type UpdateUserParams struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUser,
		arg.ID,
		arg.Email,
		arg.Phone,
		arg.Address,
		arg.UpdatedAt,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Lastname,
		&i.Firstname,
		&i.HashedPassword,
		&i.Phone,
		&i.Address,
		&i.Email,
		&i.IsAdmin,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateUserPassword = `-- name: UpdateUserPassword :one
UPDATE users SET hashed_password = $2, updated_at = $3 WHERE id = $1 RETURNING id, lastname, firstname, hashed_password, phone, address, email, is_admin, created_at, updated_at
`

type UpdateUserPasswordParams struct {
	ID             int64     `json:"id"`
	HashedPassword string    `json:"hashed_password"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUserPassword, arg.ID, arg.HashedPassword, arg.UpdatedAt)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Lastname,
		&i.Firstname,
		&i.HashedPassword,
		&i.Phone,
		&i.Address,
		&i.Email,
		&i.IsAdmin,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
