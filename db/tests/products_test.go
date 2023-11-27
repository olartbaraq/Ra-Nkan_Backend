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
	category := createRandomCategory(t)
	sub_category := createRandomSubCategory(t)

	arg := db.CreateProductParams{
		Name:          utils.RandomName(),
		Description:   utils.RandomText(),
		Price:         utils.RandomPrice(),
		Image:         "https://imagesget.com",
		QtyAval:       utils.RandomQty(),
		ShopID:        shop.ID,
		CategoryID:    category.ID,
		SubCategoryID: sub_category.ID,
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
		Name:          productTemplate.Name,
		Price:         productTemplate.Price,
		Description:   productTemplate.Description,
		Image:         productTemplate.Image,
		QtyAval:       productTemplate.QtyAval,
		ShopID:        productTemplate.ShopID,
		CategoryID:    productTemplate.CategoryID,
		SubCategoryID: productTemplate.SubCategoryID,
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
	assert.Equal(t, len(getProducts), 1)
}

func TestGetProductBySubCategory(t *testing.T) {

	product := createRandomProduct(t)

	getProducts, err := testQueries.GetProductBySubCategory(context.Background(), product.SubCategoryID)
	assert.NoError(t, err)
	assert.NotEmpty(t, getProducts)
	assert.NotEqual(t, len(getProducts), 0)
}

func TestGetProductByCategory(t *testing.T) {

	product := createRandomProduct(t)

	getProducts, err := testQueries.GetProductByCategory(context.Background(), product.CategoryID)
	assert.NoError(t, err)
	assert.NotEmpty(t, getProducts)
	assert.NotEqual(t, len(getProducts), 0)
}

func TestGetProductByPrice(t *testing.T) {

	product := createRandomProduct(t)

	getProducts, err := testQueries.GetProductByPrice(context.Background(), product.Price)
	assert.NoError(t, err)
	assert.NotEmpty(t, getProducts)
	assert.NotEqual(t, len(getProducts), 0)
}

func TestGetProductByPCS(t *testing.T) {

	product := createRandomProduct(t)

	arg := db.GetProductByPCSParams{
		Price:         product.Price,
		CategoryID:    product.CategoryID,
		SubCategoryID: product.SubCategoryID,
	}

	getProducts, err := testQueries.GetProductByPCS(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, getProducts)
	assert.NotEqual(t, len(getProducts), 0)
}

func TestListAllProducts(t *testing.T) {

	for i := 0; i < 10; i++ {
		createRandomProduct(t)
	}
	arg := db.ListAllProductParams{
		Limit:  10,
		Offset: 4,
	}

	allProducts, err := testQueries.ListAllProduct(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, allProducts)
	assert.Equal(t, int32(len(allProducts)), arg.Limit)
}

func TestUpdateProduct(t *testing.T) {

	product := createRandomProduct(t)

	arg := db.UpdateProductParams{
		ID:          product.ID,
		Name:        utils.RandomName(),
		Image:       utils.RandomEmail(),
		Price:       utils.RandomPrice(),
		Description: utils.RandomText(),
		QtyAval:     utils.RandomQty(),
		UpdatedAt:   time.Now(),
	}

	updatedProduct, err := testQueries.UpdateProduct(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, product)
	assert.Equal(t, updatedProduct.Name, arg.Name)
	assert.Equal(t, updatedProduct.Image, arg.Image)
	assert.Equal(t, updatedProduct.Price, arg.Price)
	assert.Equal(t, updatedProduct.Description, arg.Description)
	assert.Equal(t, updatedProduct.QtyAval, arg.QtyAval)
	assert.WithinDuration(t, updatedProduct.UpdatedAt, time.Now(), 2*time.Second)

}

func TestDeleteProduct(t *testing.T) {
	product := createRandomProduct(t)

	err := testQueries.DeleteProduct(context.Background(), product.ID)
	assert.NoError(t, err)

	getProduct, err := testQueries.GetProductById(context.Background(), product.ID)
	assert.Error(t, err)
	assert.Empty(t, getProduct)

}
