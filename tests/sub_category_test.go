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

func createRandomSubCategory(t *testing.T) db.SubCategory {

	category := createRandomCategory(t)

	arg := db.CreateSubCategoryParams{
		Name:         utils.RandomName(),
		CategoryID:   category.ID,
		CategoryName: category.Name,
	}

	subCategory, err := testQueries.CreateSubCategory(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, subCategory)
	assert.NotZero(t, subCategory.ID)
	assert.WithinDuration(t, subCategory.CreatedAt, time.Now(), 2*time.Second)
	assert.WithinDuration(t, subCategory.UpdatedAt, time.Now(), 2*time.Second)

	return subCategory
}

func TestCreateSubCategory(t *testing.T) {
	subCategoryPlate := createRandomSubCategory(t)

	arg := db.CreateSubCategoryParams{
		Name:         subCategoryPlate.Name,
		CategoryID:   subCategoryPlate.CategoryID,
		CategoryName: subCategoryPlate.Name,
	}

	subCategory, err := testQueries.CreateSubCategory(context.Background(), arg)
	assert.Error(t, err)
	assert.Empty(t, subCategory)

}

func TestUpdateSubCategory(t *testing.T) {

	subCategory := createRandomSubCategory(t)

	arg := db.UpdateSubCategoryParams{
		ID:        subCategory.ID,
		Name:      utils.RandomName(),
		UpdatedAt: time.Now(),
	}

	UpdatedSubCategory, err := testQueries.UpdateSubCategory(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, UpdatedSubCategory)
	assert.WithinDuration(t, UpdatedSubCategory.UpdatedAt, time.Now(), 2*time.Second)

}

func TestGetSubCategoryById(t *testing.T) {
	subCategory := createRandomSubCategory(t)

	getSubCategory, err := testQueries.GetSubCategoryById(context.Background(), subCategory.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, getSubCategory)
	assert.Equal(t, getSubCategory.ID, subCategory.ID)
	assert.Equal(t, getSubCategory.Name, subCategory.Name)
	assert.NotZero(t, getSubCategory.ID)
}

func TestGetSubCategoryByName(t *testing.T) {
	subCategory := createRandomSubCategory(t)

	getSubCategory, err := testQueries.GetSubCategoryByName(context.Background(), subCategory.Name)
	assert.NoError(t, err)
	assert.NotEmpty(t, getSubCategory)
	assert.Equal(t, getSubCategory.Name, subCategory.Name)
}

func TestGetSubCategoryByCategory(t *testing.T) {
	subCategory := createRandomSubCategory(t)

	getSubCategory, err := testQueries.GetSubCategoryByCategory(context.Background(), subCategory.CategoryID)
	assert.NoError(t, err)
	assert.NotEmpty(t, getSubCategory)
	assert.Equal(t, len(getSubCategory), 1)
}

func TestListAllSubCategory(t *testing.T) {

	for i := 0; i < 10; i++ {
		createRandomSubCategory(t)
	}
	arg := db.ListAllSubCategoryParams{
		Limit:  10,
		Offset: 0,
	}

	allSubCategories, err := testQueries.ListAllSubCategory(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, allSubCategories)
	assert.Equal(t, int32(len(allSubCategories)), arg.Limit)
}

func TestDeleteSUbCategory(t *testing.T) {
	subCategory := createRandomSubCategory(t)

	err := testQueries.DeleteSubCategory(context.Background(), subCategory.ID)
	assert.NoError(t, err)

	getSubCategory, err := testQueries.GetSubCategoryById(context.Background(), subCategory.ID)
	assert.Error(t, err)
	assert.Empty(t, getSubCategory)

}

func TestDeleteAllSUbCategory(t *testing.T) {
	subCategory := createRandomSubCategory(t)

	err := testQueries.DeleteAllSubCategories(context.Background())
	assert.NoError(t, err)

	getSubCategory, err := testQueries.GetSubCategoryById(context.Background(), subCategory.ID)
	assert.Error(t, err)
	assert.Empty(t, getSubCategory)

}
