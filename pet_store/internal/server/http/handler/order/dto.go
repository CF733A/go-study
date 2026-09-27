package order

import (
	"pet-store/internal/entities"
	"time"
)

// CreateOrderRequest represents the request structure for creating a new order
// @Description Request payload for creating a new pet store order
type CreateOrderRequest struct {
	PetID    int64                `json:"petId" example:"123"`
	Quantity int32                `json:"quantity" example:"1"`
	ShipDate time.Time            `json:"shipDate" example:"2023-12-31T15:04:05Z"`
	Status   entities.OrderStatus `json:"status" example:"approved"`
	Complete bool                 `json:"complete" example:"false"`
}

// CreateOrderResponse represents the response structure after creating an order
// @Description Response structure containing created order details
type CreateOrderResponse struct {
	ID       int64                `json:"id" example:"1"`
	PetID    int64                `json:"petId" example:"123"`
	Quantity int32                `json:"quantity" example:"1"`
	ShipDate time.Time            `json:"shipDate" example:"2023-12-31T15:04:05Z"`
	Status   entities.OrderStatus `json:"status" example:"approved"`
	Complete bool                 `json:"complete" example:"false"`
}

// GetOrderResponse represents the response structure for retrieving order details
// @Description Response structure containing complete order information
type GetOrderResponse struct {
	ID       int64                `json:"id" example:"1"`
	PetID    int64                `json:"petId" example:"123"`
	UserID   int64                `json:"userId" example:"456"`
	Quantity int32                `json:"quantity" example:"1"`
	ShipDate time.Time            `json:"shipDate" example:"2023-12-31T15:04:05Z"`
	Status   entities.OrderStatus `json:"status" example:"approved"`
	Complete bool                 `json:"complete" example:"false"`
}

// InventoryResponse represents the inventory status response
// @Description Response structure containing inventory counts by status
type InventoryResponse struct {
	Available int32 `json:"available" example:"100"`
	Pending   int32 `json:"pending" example:"15"`
	Sold      int32 `json:"sold" example:"235"`
}
