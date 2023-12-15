package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	db "github.com/olartbaraq/spectrumshelf/db/sqlc"
	"github.com/olartbaraq/spectrumshelf/utils"
)

type Product struct {
	server *Server
}

type CreateProductParams struct {
	Name            string `json:"name" binding:"required"`
	Description     string `json:"description" binding:"required"`
	Price           string `json:"price" binding:"required,numeric,isPositive"`
	Image           string `json:"image_url" binding:"required,url,isImageURL"`
	QtyAval         int32  `json:"qty_aval" binding:"required,numeric,gt=0"`
	ShopID          int64  `json:"shop_id" binding:"required"`
	ShopName        string `json:"shop_name" binding:"required"`
	CategoryID      int64  `json:"category_id" binding:"required"`
	SubCategoryID   int64  `json:"subcategory_id" binding:"required"`
	CategoryName    string `json:"category_name" binding:"required"`
	SubCategoryName string `json:"subcategory_name" binding:"required"`
}

func (p Product) router(server *Server) {

	p.server = server

	serverGroup := server.router.Group("/products", AuthenticatedMiddleware())
	serverGroup.POST("/create_product", p.createProduct)
	//serverGroup.POST("/login", a.login)
}

func (p *Product) createProduct(ctx *gin.Context) {

	incorrectImageResp := []string{
		"image URL doesn't point to a real image",
	}

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

	product := CreateProductParams{}

	if err := ctx.ShouldBindJSON(&product); err != nil {
		//fmt.Println(err.Error())
		stringErr := string(err.Error())
		//fmt.Println(stringErr)
		if strings.Contains(stringErr, "isImageURL") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Error": incorrectImageResp,
			})
			return

		}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	arg := db.CreateProductParams{
		Name:            strings.ToLower(product.Name),
		Price:           product.Price,
		Description:     product.Description,
		Image:           product.Image,
		QtyAval:         product.QtyAval,
		ShopID:          product.ShopID,
		ShopName:        strings.ToLower(product.ShopName),
		CategoryID:      product.CategoryID,
		CategoryName:    strings.ToLower(product.CategoryName),
		SubCategoryID:   product.SubCategoryID,
		SubCategoryName: strings.ToLower(product.SubCategoryName),
	}

	productToSave, err := p.server.queries.CreateProduct(context.Background(), arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "product created successfully",
		"data":    productToSave,
	})
}
