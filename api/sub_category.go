package api

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/olartbaraq/spectrumshelf/db/sqlc"
	"github.com/olartbaraq/spectrumshelf/utils"
)

type SubCategory struct {
	server *Server
}

type CreateSubCategoryParams struct {
	Name         string `json:"name" binding:"required"`
	CategoryID   int64  `json:"category_id" binding:"required"`
	CategoryName string `json:"category_name" binding:"required"`
}

type SearchSubCategoryParams struct {
	CategoryName string `json:"category_name" binding:"required"`
}

type SubCategoryResponse struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	CategoryName string    `json:"category_name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (subc SubCategory) router(server *Server) {
	subc.server = server

	serverGroup := server.router.Group("/subcategory", AuthenticatedMiddleware())
	serverGroup.POST("/create_subcategory", subc.createSubCategory)
	serverGroup.POST("/search_subcategory", subc.searchSubCategory)
	// serverGroup.GET("/list_subcategories", subc.listCategories)
}

func (subc *SubCategory) createSubCategory(ctx *gin.Context) {
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

	subCategory := CreateSubCategoryParams{}

	if err := ctx.ShouldBindJSON(&subCategory); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	arg := db.CreateSubCategoryParams{
		Name:         strings.ToLower(subCategory.Name),
		CategoryID:   subCategory.CategoryID,
		CategoryName: strings.ToLower(subCategory.CategoryName),
	}

	subCategorytoSave, err := subc.server.queries.CreateSubCategory(context.Background(), arg)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505":
				// to check for unique constraint
				stringErr := string(pqErr.Detail)
				if strings.Contains(stringErr, "name") {
					ctx.JSON(http.StatusBadRequest, gin.H{
						"Error": "SubCategory already exists",
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
		"message": "subcategory created successfully",
		"data":    subCategorytoSave,
	})
}

func (subc *SubCategory) searchSubCategory(ctx *gin.Context) {

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

	category_name := SearchSubCategoryParams{}

	if err := ctx.ShouldBindJSON(&category_name); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	allSubCategories, err := subc.server.queries.GetSubCategoryByCategory(context.Background(), strings.ToLower(category_name.CategoryName))

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
		"message": fmt.Sprintf("subcategories that belonged to %s retrieved successfully", strings.ToLower(category_name.CategoryName)),
		"data":    allSubCategories,
	})
}
