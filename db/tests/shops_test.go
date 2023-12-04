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

func createRandomShop(t *testing.T) db.Shop {
	arg := db.CreateShopParams{
		Name:    utils.RandomName(),
		Email:   utils.RandomEmail(),
		Phone:   utils.RandomPhone(),
		Address: utils.RandomAddress(),
	}

	shop, err := testQueries.CreateShop(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, shop)
	assert.Equal(t, shop.Email, arg.Email)
	assert.Equal(t, shop.Name, arg.Name)
	assert.NotZero(t, shop.ID)
	assert.WithinDuration(t, shop.CreatedAt, time.Now(), 2*time.Second)
	assert.WithinDuration(t, shop.UpdatedAt, time.Now(), 2*time.Second)

	return shop
}

func TestCreateShop(t *testing.T) {
	shopTemplate := createRandomShop(t)

	shop, err := testQueries.CreateShop(context.Background(), db.CreateShopParams{
		Name:    shopTemplate.Name,
		Phone:   shopTemplate.Phone,
		Address: shopTemplate.Address,
		Email:   shopTemplate.Email,
	})
	assert.Error(t, err)
	assert.Empty(t, shop)

}

func TestGetShopById(t *testing.T) {
	shop := createRandomShop(t)

	getShop, err := testQueries.GetShopById(context.Background(), shop.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, getShop)
	assert.Equal(t, getShop.ID, shop.ID)
	assert.Equal(t, getShop.Name, shop.Name)
}

func TestGetShopByEmail(t *testing.T) {
	shop := createRandomShop(t)

	getShop, err := testQueries.GetShopByEmail(context.Background(), shop.Email)
	assert.NoError(t, err)
	assert.NotEmpty(t, getShop)
	assert.Equal(t, getShop.ID, shop.ID)
	assert.Equal(t, getShop.Name, shop.Name)
	assert.Equal(t, getShop.Email, shop.Email)
}

func TestListAllShops(t *testing.T) {

	for i := 0; i < 10; i++ {
		createRandomShop(t)
	}
	arg := db.ListAllShopsParams{
		Limit:  10,
		Offset: 2,
	}

	allUsers, err := testQueries.ListAllShops(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, allUsers)
	assert.Equal(t, int32(len(allUsers)), arg.Limit)
}

func TestUpdateShop(t *testing.T) {

	shop := createRandomShop(t)

	arg := db.UpdateShopParams{
		ID:        shop.ID,
		Name:      utils.RandomName(),
		Email:     utils.RandomEmail(),
		Phone:     utils.RandomPhone(),
		Address:   "74 Avenue Suite, idiroko, yanibo, ajah",
		UpdatedAt: time.Now(),
	}

	updatedShop, err := testQueries.UpdateShop(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, shop)
	assert.Equal(t, updatedShop.Email, arg.Email)
	assert.Equal(t, updatedShop.Name, arg.Name)
	assert.Equal(t, updatedShop.Phone, arg.Phone)
	assert.Equal(t, updatedShop.Address, arg.Address)
	assert.WithinDuration(t, updatedShop.UpdatedAt, time.Now(), 2*time.Second)

}

func TestDeleteShop(t *testing.T) {
	shop := createRandomShop(t)

	err := testQueries.DeleteShop(context.Background(), shop.ID)
	assert.NoError(t, err)

	getShop, err := testQueries.GetShopById(context.Background(), shop.ID)
	assert.Error(t, err)
	assert.Empty(t, getShop)

}

func TestDeleteAllShop(t *testing.T) {
	shop := createRandomShop(t)

	err := testQueries.DeleteAllShops(context.Background())
	assert.NoError(t, err)

	getShop, err := testQueries.GetShopById(context.Background(), shop.ID)
	assert.Error(t, err)
	assert.Empty(t, getShop)

}
