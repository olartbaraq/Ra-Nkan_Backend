package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/olartbaraq/spectrumshelf/db/sqlc"
	"github.com/olartbaraq/spectrumshelf/utils"
)

type Shop struct {
	server *Server
}

type CreateShopParams struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
	Phone   string `json:"phone" binding:"required,len=11"`
	Address string `json:"address" binding:"required"`
}

func (s Shop) router(server *Server) {
	s.server = server
	serverGroup := server.router.Group("/shops", AuthenticatedMiddleware())
	serverGroup.POST("/create_shops", s.createShop)
	// serverGroup.GET("/allShops", s.listShops)
	// serverGroup.GET("/getShop", s.getShop)
	// serverGroup.PUT("/update_shop", s.updateShop)
}

func (s *Shop) createShop(ctx *gin.Context) {

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

	shop := CreateShopParams{}

	if err := ctx.ShouldBindJSON(&shop); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	arg := db.CreateShopParams{
		Name:    strings.ToLower(shop.Name),
		Phone:   shop.Phone,
		Address: strings.ToLower(shop.Address),
		Email:   strings.ToLower(shop.Email),
	}

	shoptoSave, err := s.server.queries.CreateShop(context.Background(), arg)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505":
				// to check for unique constraint
				stringErr := string(pqErr.Detail)
				if strings.Contains(stringErr, "name") {
					ctx.JSON(http.StatusBadRequest, gin.H{
						"Error": "shop already exists",
					})
					ctx.Abort()
					return
				} else if strings.Contains(stringErr, "phone") {
					ctx.JSON(http.StatusBadRequest, gin.H{
						"Error": "Shop with phone number already exists",
					})
					ctx.Abort()
					return
				} else if strings.Contains(stringErr, "email") {
					ctx.JSON(http.StatusBadRequest, gin.H{
						"Error": "Shop with email address already exists",
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
		"message": "shop created successfully",
		"data":    shoptoSave,
	})
}
