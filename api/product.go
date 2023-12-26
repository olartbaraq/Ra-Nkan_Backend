/*
 *   Copyright (c) 2023 Mubaraq Akanbi
 *   All rights reserved.
 *   Created by Mubaraq Akanbi
 */
package api

import (
	"context"
	"database/sql"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	db "github.com/olartbaraq/spectrumshelf/db/sqlc"
	"github.com/olartbaraq/spectrumshelf/utils"
)

type Product struct {
	server *Server
}

type File struct {
	File multipart.File `json:"image" binding:"required"`
}

type Url struct {
	Url string `json:"url" binding:"required,url"`
}

type UploadResult struct {
	Url string
	Err error
}

//	type UrlImages struct {
//		ImageUrls []string `json:"image_urls" binding:"required,isImageURL"`
//	}
type CreateProductParamsUrl struct {
	Name            string   `json:"name" binding:"required"`
	Description     string   `json:"description" binding:"required"`
	Price           string   `json:"price" binding:"required,numeric,isPositive"`
	Images          []string `json:"image_urls" binding:"required,isImageURL"`
	QtyAval         int32    `json:"qty_aval" binding:"required,numeric,gt=0"`
	ShopID          int64    `json:"shop_id" binding:"required"`
	ShopName        string   `json:"shop_name" binding:"required"`
	CategoryID      int64    `json:"category_id" binding:"required"`
	SubCategoryID   int64    `json:"subcategory_id" binding:"required"`
	CategoryName    string   `json:"category_name" binding:"required"`
	SubCategoryName string   `json:"subcategory_name" binding:"required"`
}
type CreateProductParamsFile struct {
	Name            string `json:"name" binding:"required"`
	Description     string `json:"description" binding:"required"`
	Price           string `json:"price" binding:"required,numeric,isPositive"`
	QtyAval         int32  `json:"qty_aval" binding:"required,numeric,gt=0"`
	ShopID          int64  `json:"shop_id" binding:"required"`
	ShopName        string `json:"shop_name" binding:"required"`
	CategoryID      int64  `json:"category_id" binding:"required"`
	SubCategoryID   int64  `json:"subcategory_id" binding:"required"`
	CategoryName    string `json:"category_name" binding:"required"`
	SubCategoryName string `json:"subcategory_name" binding:"required"`
}

type GetProductParams struct {
	ID   int64  `form:"id"`
	Name string `form:"name"`
}

func (p Product) router(server *Server) {

	p.server = server

	serverGroup := server.router.Group("/products")
	serverGroup.POST("/create_product_file", p.createProductByFile, AuthenticatedMiddleware())
	serverGroup.POST("/create_product_url", p.createProductByUrl, AuthenticatedMiddleware())
	serverGroup.GET("/get_products_orders", p.getProductOrders)
	serverGroup.GET("/get_product_by_id", p.getProductById)
	serverGroup.GET("/get_products_by_name", p.getProductByName)

	//serverGroup.POST("/login", a.login)
}

func NewUrlFromString(rawUrl string) (*Url, error) {
	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}

	return &Url{Url: parsedUrl.String()}, nil
}

func (p *Product) createProductByFile(ctx *gin.Context) {

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
			"status": "failed to verify token when creating products",
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

	form, err := ctx.MultipartForm()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"Error":      err.Error(),
			"data":       "Select a file to upload",
		})
		return
	}

	files := form.File["file[]"]

	uploadResults := make(chan UploadResult, len(files))

	for _, file := range files {
		go func(imageFile *multipart.FileHeader, results chan UploadResult) {
			//log.Println("WE GOT TO FILE UPLOAD HERE >>>", imageFile.Filename)

			localFile := imageFile // Create a local variable

			fileContent, err := localFile.Open()
			if err != nil {
				results <- UploadResult{Err: err}
				return
			}
			defer fileContent.Close()

			uploadUrl, err := NewMediaUpload().FileUpload(File{File: fileContent})
			// log.Println("INSIDE FILE UPLOAD GOROUTINE FILE\n", uploadUrl)
			// log.Println("INSIDE FILE UPLOAD GOROUTINE ERR", err)
			results <- UploadResult{Url: uploadUrl, Err: err}
		}(file, uploadResults)
	}

	var uploadUrls []string
	var uploadErrors []error

	for i := 0; i < len(files); i++ {
		result := <-uploadResults
		if result.Err != nil {
			uploadErrors = append(uploadErrors, result.Err)
		} else {
			uploadUrls = append(uploadUrls, result.Url)
		}
	}

	close(uploadResults)

	if len(uploadErrors) > 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"Error":      "Error uploading files",
			"data":       uploadErrors,
		})
		return
	}

	ctx.JSON(
		http.StatusPartialContent,
		gin.H{
			"message":    "files uploaded successfully to Cloudinary Server",
			"data":       uploadUrls,
			"statusCode": http.StatusPartialContent,
		},
	)

	product := CreateProductParamsFile{}

	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	arg := db.CreateProductParams{
		Name:            strings.ToLower(product.Name),
		Price:           product.Price,
		Description:     product.Description,
		Images:          uploadUrls,
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

func (p *Product) createProductByUrl(ctx *gin.Context) {

	imageValidationResponse := []string{
		"one or more image URL doesn't point to a real image",
		"one or more image URL size is more than 500kb",
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

	product := CreateProductParamsUrl{}

	if err := ctx.ShouldBindJSON(&product); err != nil {
		//fmt.Println(err.Error())
		stringErr := string(err.Error())
		//fmt.Println(stringErr)
		if strings.Contains(stringErr, "isImageURL") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"Error": imageValidationResponse,
			})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	uploadResults := make(chan UploadResult, len(product.Images))

	for _, imageUrl := range product.Images {
		go func(imageUrl string, results chan<- UploadResult) {
			//log.Println("WE GOT to REMOTE UPLOAD HERE >>>", imageUrl)

			parsedUrl, err := NewUrlFromString(imageUrl)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"Error": err.Error(),
				})
				return
			}

			uploadUrl, err := NewMediaUpload().RemoteUpload(*parsedUrl)
			results <- UploadResult{Url: uploadUrl, Err: err}
		}(imageUrl, uploadResults)
	}

	var uploadUrls []string
	var uploadErrors []error

	for i := 0; i < len(product.Images); i++ {
		result := <-uploadResults
		if result.Err != nil {
			uploadErrors = append(uploadErrors, result.Err)
		} else {
			uploadUrls = append(uploadUrls, result.Url)
		}
	}

	close(uploadResults)

	if len(uploadErrors) > 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"Error":      "Error uploading files",
			"data":       uploadErrors,
		})
		return
	}

	// ctx.JSON(
	// 	http.StatusOK,
	// 	gin.H{
	// 		"statusCode": http.StatusOK,
	// 		"message":    "files uploaded successfully",
	// 		"data":       uploadUrls,
	// 	},
	// )

	arg := db.CreateProductParams{
		Name:            strings.ToLower(product.Name),
		Price:           product.Price,
		Description:     product.Description,
		Images:          uploadUrls,
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

func (p *Product) getProductOrders(ctx *gin.Context) {

	productByOrders, err := p.server.queries.ListAllProductsByOrders(context.Background())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "products retrieved successfully",
		"data":    productByOrders,
	})
}

func (p *Product) getProductById(ctx *gin.Context) {

	product := GetProductParams{}

	if err := ctx.ShouldBindQuery(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	productByID, err := p.server.queries.GetProductById(context.Background(), product.ID)
	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{
			"Error":   err.Error(),
			"message": "Product not found",
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
		"message": "product retrieved successfully",
		"data":    productByID,
	})
}

func (p *Product) getProductByName(ctx *gin.Context) {

	product := GetProductParams{}

	if err := ctx.ShouldBindQuery(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	if strings.TrimSpace(product.Name) == "" {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "no input entered",
		})
		return
	}

	nullString := sql.NullString{String: product.Name, Valid: true}

	productsByName, err := p.server.queries.GetProductByName(context.Background(), nullString)

	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error":   err.Error(),
			"message": "Issue Encountered, try again later",
		})
		ctx.Abort()
		return
	}

	if len(productsByName) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Product not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "products retrieved successfully",
		"data":    productsByName,
	})
}
