package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/olartbaraq/spectrumshelf/db/sqlc"
	"github.com/olartbaraq/spectrumshelf/utils"
)

type Order struct {
	server *Server
}

type OrderItem struct {
	ProductID  int     `json:"product_id"`
	Name       string  `json:"name"`
	Image      string  `json:"image"`
	ShopName   string  `json:"shop_name"`
	QtyBought  int     `json:"qty_bought"`
	UnitPrice  float64 `json:"unit_price"`
	TotalPrice float64 `json:"total_price"`
}

type CreateOrderParams struct {
	UserID         int64       `json:"user_id" binding:"required"`
	Items          []OrderItem `json:"items" binding:"required"`
	TransactionRef string      `json:"transaction_ref"`
	TotalPrice     string      `json:"total_price"`
	Status         string      `json:"status"`
}
type CompleteOrderParams struct {
	ID     int64  `json:"id" binding:"required"`
	PayRef string `json:"pay_ref" binding:"required"`
}

type OrderResponse struct {
	ID        int64       `json:"id"`
	UserID    int64       `json:"user_id"`
	Items     []OrderItem `json:"items"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

func (order Order) router(server *Server) {
	order.server = server

	serverGroup := server.router.Group("/order", AuthenticatedMiddleware())

	serverGroup.POST("/create_order", order.createOrder)
	serverGroup.PUT("/complete_order", order.completeOrder)
}

func convertRawMessageToOrderItems(rawMessage json.RawMessage) ([]OrderItem, error) {
	var orderItems []OrderItem
	if err := json.Unmarshal(rawMessage, &orderItems); err != nil {
		return nil, err
	}
	return orderItems, nil
}

var randomTransRef string

func (o *Order) createOrder(ctx *gin.Context) {

	tokenString, err := extractTokenFromRequest(ctx)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: Missing or invalid token",
		})
		return
	}

	userId, _, err := returnIdRole(tokenString)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error":  err.Error(),
			"status": "failed to verify token",
		})
		ctx.Abort()
		return
	}

	order := CreateOrderParams{}

	if err := ctx.ShouldBindJSON(&order); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	if userId != order.UserID {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: invalid token",
		})
		ctx.Abort()
		return
	}

	itemsJSON, err := json.Marshal(order.Items)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
		return
	}

	randomTransRef = fmt.Sprintf("%v-%v", utils.RandIntegers(3), utils.RandomString(10))

	wg := sync.WaitGroup{}
	newCtx, cancel := context.WithCancel(ctx)
	if err := newCtx.Err(); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"Error":   err.Error(),
			"message": "error creating context",
		})
	}

	defer cancel()

	var overallPrice float64
	var totalPriceChan = make(chan float64)

	for _, value := range order.Items {

		wg.Add(1)

		go func(eachItem OrderItem, priceChan chan float64) {

			defer wg.Done()

			totalPrice := eachItem.TotalPrice

			priceChan <- totalPrice

		}(value, totalPriceChan)
	}

	go func() {
		wg.Wait()
		close(totalPriceChan)
	}()

	for totalPrice := range totalPriceChan {
		overallPrice += totalPrice
	}

	overallPriceString := fmt.Sprintf("%v", overallPrice)

	arg := db.CreateOrderParams{
		UserID:         order.UserID,
		Items:          json.RawMessage(itemsJSON),
		TransactionRef: randomTransRef,
		PayRef:         utils.RandomString(1),
		TotalPrice:     overallPriceString,
		Status:         "pending",
	}

	newOrder, err := o.server.queries.CreateOrder(context.Background(), arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "order created successfully",
		"data":    newOrder,
	})
}

func (o *Order) completeOrder(ctx *gin.Context) {

	order := CompleteOrderParams{}

	if err := ctx.ShouldBindJSON(&order); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	arg := db.CompleteOrderParams{
		ID:     order.ID,
		PayRef: order.PayRef,
		Status: "completed",
	}

	newOrder, err := o.server.queries.CompleteOrder(context.Background(), arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
		return
	}

	orderItems, err := convertRawMessageToOrderItems(newOrder.Items)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
		return
	}

	wg := sync.WaitGroup{}
	newCtx, cancel := context.WithCancel(ctx)
	if err := newCtx.Err(); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"Error":   err.Error(),
			"message": "error creating context",
		})
	}

	defer cancel()

	for _, value := range orderItems {
		wg.Add(1)

		go func(item OrderItem) {

			defer wg.Done()

			productCtx, productCancel := context.WithCancel(newCtx)
			defer productCancel()

			select {
			case <-productCtx.Done():
				ctx.JSON(http.StatusNotFound, gin.H{
					"Error":   productCtx.Err().Error(),
					"message": "error creating context inside goroutine",
				})
				ctx.Abort()
				return

			default:
				productByID, err := o.server.queries.GetProductById(productCtx, int64(item.ProductID))
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

				arg := db.UpdateProductParams{
					ID:          productByID.ID,
					Name:        productByID.Name,
					Images:      productByID.Images,
					Price:       productByID.Price,
					Description: productByID.Description,
					QtyAval:     productByID.QtyAval - int32(item.QtyBought),
					UpdatedAt:   time.Now(),
				}

				_, err = o.server.queries.UpdateProduct(context.Background(), arg)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"Error":   err.Error(),
						"message": "Issue Encountered updating product, try again later",
					})
					return
				}
			}
		}(value)
	}

	wg.Wait()

	orderResponse := OrderResponse{
		ID:        newOrder.ID,
		UserID:    newOrder.UserID,
		Items:     orderItems,
		CreatedAt: newOrder.CreatedAt,
		UpdatedAt: newOrder.UpdatedAt,
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "order completed successfully",
		"data":    orderResponse,
	})
}
