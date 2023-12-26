package api

import (
	"context"
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/olartbaraq/spectrumshelf/db/sqlc"
	"github.com/olartbaraq/spectrumshelf/utils"
)

type Category struct {
	server *Server
}

type CreateCategoryParams struct {
	Name string `json:"name" binding:"required"`
}

type UpdateCategoryParams struct {
	ID   int64  `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type CategoryResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c Category) router(server *Server) {
	c.server = server

	serverGroup := server.router.Group("/category", AuthenticatedMiddleware())
	serverGroup.POST("/create_category", c.createCategory)
	serverGroup.POST("/search_category", c.searchCategory)
	serverGroup.GET("/list_categories", c.listCategories)
	serverGroup.PUT("/update_category", c.updateCategory)
}

func (c *Category) createCategory(ctx *gin.Context) {

	tokenString, err := extractTokenFromRequest(ctx)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: Missing or invalid token",
		})
		return
	}

	_, role, err := returnIdRole(tokenString)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error":  err.Error(),
			"status": "failed to verify token",
		})
		ctx.Abort()
		return
	}

	if role != utils.AdminRole {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		ctx.Abort()
		return
	}

	category := CreateCategoryParams{}

	if err := ctx.ShouldBindJSON(&category); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	categorytoSave, err := c.server.queries.CreateCategory(context.Background(), strings.ToLower(category.Name))

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505":
				// to check for unique constraint
				stringErr := string(pqErr.Detail)
				if strings.Contains(stringErr, "name") {
					ctx.JSON(http.StatusBadRequest, gin.H{
						"Error": "Category already exists",
					})
					ctx.Abort()
					return
				}
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"Error": err.Error(),
				})
				ctx.Abort()
				return
			}
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Error": err.Error(),
			})
			ctx.Abort()
			return
		}
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "category created successfully",
		"data":    categorytoSave,
	})
}

func (c *Category) searchCategory(ctx *gin.Context) {

	tokenString, err := extractTokenFromRequest(ctx)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: Missing or invalid token",
		})
		return
	}

	_, role, err := returnIdRole(tokenString)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error":  err.Error(),
			"status": "failed to verify token",
		})
		ctx.Abort()
		return
	}

	if role != utils.AdminRole {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		ctx.Abort()
		return
	}

	category := CreateCategoryParams{}

	if err := ctx.ShouldBindJSON(&category); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	getCategory, err := c.server.queries.GetCategoryByName(context.Background(), strings.ToLower(category.Name))

	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{
			"Error":   err.Error(),
			"message": "Category not found",
		})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error":   err.Error(),
			"message": "Issue Encountered, try again later",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "category retrieved successfully",
		"data":    getCategory,
	})
}

func (c *Category) listCategories(ctx *gin.Context) {

	tokenString, err := extractTokenFromRequest(ctx)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: Missing or invalid token",
		})
		return
	}

	_, role, err := returnIdRole(tokenString)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error":  err.Error(),
			"status": "failed to verify token",
		})
		ctx.Abort()
		return
	}

	if role != utils.AdminRole {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		ctx.Abort()
		return
	}

	arg := db.ListAllCategoryParams{
		Limit:  10,
		Offset: 0,
	}

	categories, err := c.server.queries.ListAllCategory(context.Background(), arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
		return
	}

	allCategories := []CategoryResponse{} // slice of type struct

	for _, v := range categories {

		categoryResponse := CategoryResponse{
			ID:        v.ID,
			Name:      v.Name,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		}
		n := categoryResponse
		allCategories = append(allCategories, n)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "all categories fetched sucessfully",
		"data":    allCategories,
	})
}

func (c *Category) updateCategory(ctx *gin.Context) {

	tokenString, err := extractTokenFromRequest(ctx)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: Missing or invalid token",
		})
		return
	}

	_, role, err := returnIdRole(tokenString)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error":  err.Error(),
			"status": "failed to verify token",
		})
		ctx.Abort()
		return
	}

	if role != utils.AdminRole {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		ctx.Abort()
		return
	}

	category := UpdateCategoryParams{}

	if err := ctx.ShouldBindJSON(&category); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	arg := db.UpdateCategoryParams{
		ID:        category.ID,
		Name:      category.Name,
		UpdatedAt: time.Now(),
	}

	categoryToUpdate, err := c.server.queries.UpdateCategory(context.Background(), arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "category updated successfully",
		"data":    categoryToUpdate,
	})
}
