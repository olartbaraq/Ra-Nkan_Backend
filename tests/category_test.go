package all_test

import (
	"context"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"

	db "github.com/olartbaraq/spectrumshelf/db/sqlc"
	"github.com/olartbaraq/spectrumshelf/utils"
)

func createRandomCategory(t *testing.T) db.Category {

	category, err := testQueries.CreateCategory(context.Background(), utils.RandomName())
	assert.NoError(t, err)
	assert.NotEmpty(t, category)
	assert.NotZero(t, category.ID)
	assert.WithinDuration(t, category.CreatedAt, time.Now(), 2*time.Second)
	assert.WithinDuration(t, category.UpdatedAt, time.Now(), 2*time.Second)

	return category
}

func TestCreateCategory(t *testing.T) {
	categoryPlate := createRandomCategory(t)

	category, err := testQueries.CreateCategory(context.Background(), categoryPlate.Name)
	assert.Error(t, err)
	assert.Empty(t, category)

}

func TestUpdateCategory(t *testing.T) {

	category := createRandomCategory(t)

	arg := db.UpdateCategoryParams{
		ID:        category.ID,
		Name:      utils.RandomName(),
		UpdatedAt: time.Now(),
	}

	UpdatedCategory, err := testQueries.UpdateCategory(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, UpdatedCategory)
	assert.WithinDuration(t, UpdatedCategory.UpdatedAt, time.Now(), 2*time.Second)

}

func TestGetCategoryById(t *testing.T) {
	category := createRandomCategory(t)

	getCategory, err := testQueries.GetCategoryById(context.Background(), category.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, getCategory)
	assert.Equal(t, getCategory.ID, category.ID)
	assert.Equal(t, getCategory.Name, category.Name)
	assert.NotZero(t, getCategory.ID)
}

func TestGetCategoryByName(t *testing.T) {
	category := createRandomCategory(t)

	getCategory, err := testQueries.GetCategoryByName(context.Background(), category.Name)
	assert.NoError(t, err)
	assert.NotEmpty(t, getCategory)
	assert.Equal(t, getCategory.Name, category.Name)
}

func TestListAllCategory(t *testing.T) {

	for i := 0; i < 10; i++ {
		createRandomCategory(t)
	}
	arg := db.ListAllCategoryParams{
		Limit:  10,
		Offset: 0,
	}

	allCategories, err := testQueries.ListAllCategory(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, allCategories)
	assert.Equal(t, int32(len(allCategories)), arg.Limit)
}

func TestDeleteCategory(t *testing.T) {
	category := createRandomCategory(t)

	err := testQueries.DeleteCategory(context.Background(), category.ID)
	assert.NoError(t, err)

	getCategory, err := testQueries.GetCategoryById(context.Background(), category.ID)
	assert.Error(t, err)
	assert.Empty(t, getCategory)

}

func TestDeleteAllCategory(t *testing.T) {
	category := createRandomCategory(t)

	err := testQueries.DeleteAllCategories(context.Background())
	assert.NoError(t, err)

	getCategory, err := testQueries.GetCategoryById(context.Background(), category.ID)
	assert.Error(t, err)
	assert.Empty(t, getCategory)

}
