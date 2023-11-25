// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: products.sql

package db

import (
	"context"
)

const createProduct = `-- name: CreateProduct :one
INSERT INTO products (
    name,
    description,
    price,
    image,
    qty_aval,
    shop_id
) VALUES (
    $1, $2, $3, $4, $5, $6) RETURNING id, name, description, price, image, qty_aval, shop_id, created_at, updated_at
`

type CreateProductParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       string `json:"price"`
	Image       string `json:"image"`
	QtyAval     int32  `json:"qty_aval"`
	ShopID      int64  `json:"shop_id"`
}

func (q *Queries) CreateProduct(ctx context.Context, arg CreateProductParams) (Product, error) {
	row := q.db.QueryRowContext(ctx, createProduct,
		arg.Name,
		arg.Description,
		arg.Price,
		arg.Image,
		arg.QtyAval,
		arg.ShopID,
	)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Price,
		&i.Image,
		&i.QtyAval,
		&i.ShopID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteProduct = `-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1
`

func (q *Queries) DeleteProduct(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteProduct, id)
	return err
}

const getProductById = `-- name: GetProductById :one
SELECT id, name, description, price, image, qty_aval, shop_id, created_at, updated_at FROM products WHERE id = $1
`

func (q *Queries) GetProductById(ctx context.Context, id int64) (Product, error) {
	row := q.db.QueryRowContext(ctx, getProductById, id)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Price,
		&i.Image,
		&i.QtyAval,
		&i.ShopID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getProductByName = `-- name: GetProductByName :many
SELECT id, name, description, price, image, qty_aval, shop_id, created_at, updated_at FROM products WHERE name = $1 ORDER BY id
`

func (q *Queries) GetProductByName(ctx context.Context, name string) ([]Product, error) {
	rows, err := q.db.QueryContext(ctx, getProductByName, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Product{}
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Price,
			&i.Image,
			&i.QtyAval,
			&i.ShopID,
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

const getProductByShop = `-- name: GetProductByShop :many
SELECT id, name, description, price, image, qty_aval, shop_id, created_at, updated_at FROM products WHERE shop_id = $1 ORDER BY id
`

func (q *Queries) GetProductByShop(ctx context.Context, shopID int64) ([]Product, error) {
	rows, err := q.db.QueryContext(ctx, getProductByShop, shopID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Product{}
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Price,
			&i.Image,
			&i.QtyAval,
			&i.ShopID,
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

const listAllProduct = `-- name: ListAllProduct :many
SELECT id, name, description, price, image, qty_aval, shop_id, created_at, updated_at FROM products ORDER BY id LIMIT $1 OFFSET $2
`

type ListAllProductParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListAllProduct(ctx context.Context, arg ListAllProductParams) ([]Product, error) {
	rows, err := q.db.QueryContext(ctx, listAllProduct, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Product{}
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Price,
			&i.Image,
			&i.QtyAval,
			&i.ShopID,
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

const updateProduct = `-- name: UpdateProduct :one
UPDATE products SET name = $2, qty_aval = $6, description = $5, price = $4, image = $3, updated_at = $6 WHERE id = $1 RETURNING id, name, description, price, image, qty_aval, shop_id, created_at, updated_at
`

type UpdateProductParams struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Image       string `json:"image"`
	Price       string `json:"price"`
	Description string `json:"description"`
	QtyAval     int32  `json:"qty_aval"`
}

func (q *Queries) UpdateProduct(ctx context.Context, arg UpdateProductParams) (Product, error) {
	row := q.db.QueryRowContext(ctx, updateProduct,
		arg.ID,
		arg.Name,
		arg.Image,
		arg.Price,
		arg.Description,
		arg.QtyAval,
	)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Price,
		&i.Image,
		&i.QtyAval,
		&i.ShopID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
