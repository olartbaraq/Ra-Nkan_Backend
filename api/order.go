package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/olartbaraq/spectrumshelf/db/sqlc"
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
	UserID int64       `json:"user_id" binding:"required"`
	Items  []OrderItem `json:"items" binding:"required"`
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
}

func convertRawMessageToOrderItems(rawMessage json.RawMessage) ([]OrderItem, error) {
	var orderItems []OrderItem
	if err := json.Unmarshal(rawMessage, &orderItems); err != nil {
		return nil, err
	}
	return orderItems, nil
}

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

	arg := db.CreateOrderParams{
		UserID: order.UserID,
		Items:  json.RawMessage(itemsJSON),
	}

	newOrder, err := o.server.queries.CreateOrder(context.Background(), arg)

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

	orderResponse := OrderResponse{
		ID:        newOrder.ID,
		UserID:    newOrder.UserID,
		Items:     orderItems,
		CreatedAt: newOrder.CreatedAt,
		UpdatedAt: newOrder.UpdatedAt,
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"status":  "success",
		"message": "order created successfully",
		"data":    orderResponse,
	})
}
