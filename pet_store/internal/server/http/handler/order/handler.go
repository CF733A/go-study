package order

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"pet-store/internal/entities"
	"pet-store/internal/server/http/middleware"
	"pet-store/internal/server/http/response"
	services "pet-store/internal/services"

	"github.com/go-chi/chi/v5"
)

type OrderHandler struct {
	orderService services.OrderService
}

func NewOrderHandler(orderService services.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// CreateOrder godoc
// @Summary Place an order for a pet
// @Description Place a new order in the store
// @Tags store
// @Accept json
// @Produce json
// @Param request body CreateOrderRequest true "Order data"
// @Success 201 {object} CreateOrderResponse "Order created successfully"
// @Failure 400 {object} response.Response "Invalid input data"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "Pet not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /store/order [post]
// @Example request
// {
//   "petId": 1,
//   "quantity": 1,
//   "shipDate": "2023-01-01T00:00:00Z",
//   "status": "placed",
//   "complete": false
// }
// @Example response
// {
//   "status": true,
//   "data": {
//     "id": 1,
//     "petId": 1,
//     "quantity": 1,
//     "shipDate": "2023-01-01T00:00:00Z",
//     "status": "placed",
//     "complete": false
//   }
// }
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.JSONer(w, http.StatusMethodNotAllowed, nil, errors.New("method not allowed"))
		return
	}

	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSONer(w, http.StatusBadRequest, nil, errors.New("invalid JSON format"))
		return
	}

	if req.PetID <= 0 {
		response.JSONer(w, http.StatusBadRequest, nil, errors.New("pet ID is required and must be positive"))
		return
	}

	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		response.JSONer(w, http.StatusUnauthorized, nil, errors.New("user authentication required"))
		return
	}

	order := &entities.Order{
		PetID:    req.PetID,
		UserID:   userID,
		Quantity: req.Quantity,
		ShipDate: req.ShipDate,
		Status:   req.Status,
		Complete: req.Complete,
	}

	if order.Quantity == 0 {
		order.Quantity = 1
	}

	if order.ShipDate.IsZero() {
		order.ShipDate = time.Now()
	}

	if order.Status == "" {
		order.Status = entities.StatusPlaced
	}

	if err := h.orderService.CreateOrder(r.Context(), order); err != nil {
		switch err.Error() {
		case "pet not found":
			response.JSONer(w, http.StatusNotFound, nil, err)
		case "pet is not available for purchase":
			response.JSONer(w, http.StatusBadRequest, nil, err)
		default:
			response.JSONer(w, http.StatusInternalServerError, nil, err)
		}
		return
	}

	responseData := CreateOrderResponse{
		ID:       order.ID,
		PetID:    order.PetID,
		Quantity: order.Quantity,
		ShipDate: order.ShipDate,
		Status:   order.Status,
		Complete: order.Complete,
	}

	response.JSONer(w, http.StatusCreated, responseData, nil)
}

// GetOrderByID godoc
// @Summary Find purchase order by ID
// @Description Get order details by order ID
// @Tags store
// @Produce json
// @Param orderId path int64 true "Order ID"
// @Success 200 {object} GetOrderResponse "Order details"
// @Failure 400 {object} response.Response "Invalid order ID"
// @Failure 404 {object} response.Response "Order not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /store/order/{orderId} [get]
// @Example response
// {
//   "status": true,
//   "data": {
//     "id": 1,
//     "petId": 1,
//     "userId": 1,
//     "quantity": 1,
//     "shipDate": "2023-01-01T00:00:00Z",
//     "status": "placed",
//     "complete": false
//   }
// }
func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.JSONer(w, http.StatusMethodNotAllowed, nil, errors.New("method not allowed"))
		return
	}

	orderIDStr := chi.URLParam(r, "orderId")
	if orderIDStr == "" {
		response.JSONer(w, http.StatusBadRequest, nil, errors.New("order ID is required"))
		return
	}

	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil || orderID <= 0 {
		response.JSONer(w, http.StatusBadRequest, nil, errors.New("invalid order ID"))
		return
	}

	order, err := h.orderService.GetOrderByID(r.Context(), orderID)
	if err != nil {
		if err.Error() == "order not found" {
			response.JSONer(w, http.StatusNotFound, nil, err)
		} else {
			response.JSONer(w, http.StatusInternalServerError, nil, err)
		}
		return
	}

	responseData := GetOrderResponse{
		ID:       order.ID,
		PetID:    order.PetID,
		UserID:   order.UserID,
		Quantity: order.Quantity,
		ShipDate: order.ShipDate,
		Status:   order.Status,
		Complete: order.Complete,
	}

	response.JSONer(w, http.StatusOK, responseData, nil)
}

// DeleteOrder godoc
// @Summary Delete purchase order by ID
// @Description Delete order by order ID
// @Tags store
// @Produce json
// @Param orderId path int64 true "Order ID"
// @Success 200 {object} response.Response "Order deleted successfully"
// @Failure 400 {object} response.Response "Invalid order ID"
// @Failure 404 {object} response.Response "Order not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /store/order/{orderId} [delete]
// @Example response
// {
//   "status": true,
//   "message": "order deleted successfully"
// }
func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		response.JSONer(w, http.StatusMethodNotAllowed, nil, errors.New("method not allowed"))
		return
	}

	orderIDStr := chi.URLParam(r, "orderId")
	if orderIDStr == "" {
		response.JSONer(w, http.StatusBadRequest, nil, errors.New("order ID is required"))
		return
	}

	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil || orderID <= 0 {
		response.JSONer(w, http.StatusBadRequest, nil, errors.New("invalid order ID"))
		return
	}

	if err := h.orderService.DeleteOrder(r.Context(), orderID); err != nil {
		switch err.Error() {
		case "order not found":
			response.JSONer(w, http.StatusNotFound, nil, err)
		case "cannot delete delivered order":
			response.JSONer(w, http.StatusBadRequest, nil, err)
		default:
			response.JSONer(w, http.StatusInternalServerError, nil, err)
		}
		return
	}

	response.JSONer(w, http.StatusOK, map[string]string{"message": "order deleted successfully"}, nil)
}

// GetInventory godoc
// @Summary Returns pet inventories by status
// @Description Get pet inventory counts by status
// @Tags store
// @Produce json
// @Success 200 {object} InventoryResponse "Inventory counts"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /store/inventory [get]
// @Example response
// {
//   "status": true,
//   "data": {
//     "available": 5,
//     "pending": 2,
//     "sold": 10
//   }
// }
func (h *OrderHandler) GetInventory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.JSONer(w, http.StatusMethodNotAllowed, nil, errors.New("method not allowed"))
		return
	}

	inventory, err := h.orderService.GetInventory(r.Context())
	if err != nil {
		response.JSONer(w, http.StatusInternalServerError, nil, err)
		return
	}

	responseData := InventoryResponse{
		Available: inventory["available"],
		Pending:   inventory["pending"],
		Sold:      inventory["sold"],
	}

	response.JSONer(w, http.StatusOK, responseData, nil)
}
