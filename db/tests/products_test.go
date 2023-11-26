package db_test

import (
	"context"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"

	db "github.com/olartbaraq/spectrumshelf/db/sqlc"
	"github.com/olartbaraq/spectrumshelf/utils"
)

func createRandomProduct(t *testing.T) db.Product {
	shop := createRandomShop(t)

	arg := db.CreateProductParams{
		Name:        utils.RandomName(),
		Description: utils.RandomText(),
		Price:       utils.RandomPrice(),
		Image:       "https://imagesget.com",
		QtyAval:     utils.RandomQty(),
		ShopID:      shop.ID,
	}

	product, err := testQueries.CreateProduct(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, product)
	assert.Equal(t, product.Name, arg.Name)
	assert.Equal(t, product.Price, arg.Price)
	assert.Equal(t, product.QtyAval, arg.QtyAval)
	assert.NotZero(t, product.ID)
	assert.WithinDuration(t, product.CreatedAt, time.Now(), 2*time.Second)
	assert.WithinDuration(t, product.UpdatedAt, time.Now(), 2*time.Second)

	return product
}

func TestCreateProduct(t *testing.T) {
	productTemplate := createRandomProduct(t)

	product, err := testQueries.CreateProduct(context.Background(), db.CreateProductParams{
		Name:        productTemplate.Name,
		Price:       productTemplate.Price,
		Description: productTemplate.Description,
		Image:       productTemplate.Image,
		QtyAval:     productTemplate.QtyAval,
		ShopID:      productTemplate.ShopID,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, product)

}

func TestGetProductById(t *testing.T) {
	product := createRandomProduct(t)

	getProduct, err := testQueries.GetProductById(context.Background(), product.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, getProduct)
	assert.Equal(t, getProduct.ID, product.ID)
	assert.Equal(t, getProduct.Name, product.Name)
	assert.Equal(t, getProduct.Price, product.Price)
	assert.Equal(t, getProduct.ShopID, product.ShopID)

}

func TestGetProductByName(t *testing.T) {
	product := createRandomProduct(t)

	getProducts, err := testQueries.GetProductByName(context.Background(), product.Name)
	assert.NoError(t, err)
	assert.NotEmpty(t, getProducts)
}

func TestGetProductByShop(t *testing.T) {
	product := createRandomProduct(t)

	getProducts, err := testQueries.GetProductByShop(context.Background(), product.ShopID)
	assert.NoError(t, err)
	assert.NotEmpty(t, getProducts)
}

// func TestListAllShops(t *testing.T) {
// 	arg := db.ListAllShopsParams{
// 		Limit:  10,
// 		Offset: 2,
// 	}

// 	allUsers, err := testQueries.ListAllShops(context.Background(), arg)
// 	assert.NoError(t, err)
// 	assert.NotEmpty(t, allUsers)
// }

// func TestUpdateShop(t *testing.T) {

// 	shop := createRandomShop(t)

// 	arg := db.UpdateShopParams{
// 		ID:        shop.ID,
// 		Name:      utils.RandomName(),
// 		Email:     utils.RandomEmail(),
// 		Phone:     utils.RandomPhone(),
// 		Address:   "74 Avenue Suite, idiroko, yanibo, ajah",
// 		UpdatedAt: time.Now(),
// 	}

// 	updatedShop, err := testQueries.UpdateShop(context.Background(), arg)
// 	assert.NoError(t, err)
// 	assert.NotEmpty(t, shop)
// 	assert.Equal(t, updatedShop.Email, arg.Email)
// 	assert.Equal(t, updatedShop.Name, arg.Name)
// 	assert.Equal(t, updatedShop.Phone, arg.Phone)
// 	assert.Equal(t, updatedShop.Address, arg.Address)
// 	assert.WithinDuration(t, updatedShop.UpdatedAt, time.Now(), 2*time.Second)

// }

// func TestDeleteShop(t *testing.T) {
// 	shop := createRandomShop(t)

// 	err := testQueries.DeleteShop(context.Background(), shop.ID)
// 	assert.NoError(t, err)

// }
