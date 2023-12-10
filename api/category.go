package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/olartbaraq/spectrumshelf/utils"
)

type Category struct {
	server *Server
}

type CreateCategoryParams struct {
	Name string `json:"name" binding:"required"`
}

func (c Category) router(server *Server) {
	c.server = server

	serverGroup := server.router.Group("/category", AuthenticatedMiddleware())
	serverGroup.POST("/create_category", c.createCategory)
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

	categorytoSave, err := c.server.queries.CreateCategory(context.Background(), category.Name)

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
