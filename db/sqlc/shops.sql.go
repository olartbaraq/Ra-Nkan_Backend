// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: shops.sql

package db

import (
	"context"
	"time"
)

const createShop = `-- name: CreateShop :one
INSERT INTO shops (
    name,
    phone,
    address,
    email
) VALUES (
    $1, $2, $3, $4) RETURNING id, name, phone, address, email, created_at, updated_at
`

type CreateShopParams struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	Email   string `json:"email"`
}

func (q *Queries) CreateShop(ctx context.Context, arg CreateShopParams) (Shop, error) {
	row := q.db.QueryRowContext(ctx, createShop,
		arg.Name,
		arg.Phone,
		arg.Address,
		arg.Email,
	)
	var i Shop
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Phone,
		&i.Address,
		&i.Email,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteShop = `-- name: DeleteShop :exec
DELETE FROM shops WHERE id = $1
`

func (q *Queries) DeleteShop(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteShop, id)
	return err
}

const getShopByEmail = `-- name: GetShopByEmail :one
SELECT id, name, phone, address, email, created_at, updated_at FROM shops WHERE email = $1
`

func (q *Queries) GetShopByEmail(ctx context.Context, email string) (Shop, error) {
	row := q.db.QueryRowContext(ctx, getShopByEmail, email)
	var i Shop
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Phone,
		&i.Address,
		&i.Email,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getShopById = `-- name: GetShopById :one
SELECT id, name, phone, address, email, created_at, updated_at FROM shops WHERE id = $1
`

func (q *Queries) GetShopById(ctx context.Context, id int64) (Shop, error) {
	row := q.db.QueryRowContext(ctx, getShopById, id)
	var i Shop
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Phone,
		&i.Address,
		&i.Email,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listAllShops = `-- name: ListAllShops :many
SELECT id, name, phone, address, email, created_at, updated_at FROM shops ORDER BY id LIMIT $1 OFFSET $2
`

type ListAllShopsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListAllShops(ctx context.Context, arg ListAllShopsParams) ([]Shop, error) {
	rows, err := q.db.QueryContext(ctx, listAllShops, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Shop{}
	for rows.Next() {
		var i Shop
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Phone,
			&i.Address,
			&i.Email,
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

const updateShop = `-- name: UpdateShop :one
UPDATE shops SET name = $2, address = $5, phone = $4, email = $3, updated_at = $6 WHERE id = $1 RETURNING id, name, phone, address, email, created_at, updated_at
`

type UpdateShopParams struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (q *Queries) UpdateShop(ctx context.Context, arg UpdateShopParams) (Shop, error) {
	row := q.db.QueryRowContext(ctx, updateShop,
		arg.ID,
		arg.Name,
		arg.Email,
		arg.Phone,
		arg.Address,
		arg.UpdatedAt,
	)
	var i Shop
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Phone,
		&i.Address,
		&i.Email,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
