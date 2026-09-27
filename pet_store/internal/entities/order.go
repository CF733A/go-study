package entities

import (
	"errors"
	"time"
)

type OrderStatus string

const (
	StatusPlaced    OrderStatus = "placed"
	StatusApproved  OrderStatus = "approved"
	StatusDelivered OrderStatus = "delivered"
)

type Order struct {
	ID       int64       `json:"id"`
	PetID    int64       `json:"petId"`
	UserID   int64       `json:"userId"`
	Quantity int32       `json:"quantity"`
	ShipDate time.Time   `json:"shipDate"`
	Status   OrderStatus `json:"status"`
	Complete bool        `json:"complete"`
}

func NewOrder(petID, userID int64) *Order {
	return &Order{
		PetID:    petID,
		UserID:   userID,
		Quantity: 1,
		ShipDate: time.Now(),
		Status:   StatusPlaced,
		Complete: false,
	}
}

func (o *Order) MarkAsPlaced() {
	o.Status = StatusPlaced
	o.Complete = false
}

func (o *Order) MarkAsApproved() {
	o.Status = StatusApproved
	o.Complete = false
}

func (o *Order) MarkAsDelivered() {
	o.Status = StatusDelivered
	o.Complete = true
}

func (o *Order) IsComplete() bool {
	return o.Complete
}

func (o *Order) Update(updatedOrder *Order) {
	o.PetID = updatedOrder.PetID
	o.UserID = updatedOrder.UserID
	o.ShipDate = updatedOrder.ShipDate
	o.Status = updatedOrder.Status
	o.Complete = updatedOrder.Complete
}

func (o *Order) Validate() error {
	if o.PetID <= 0 {
		return errors.New("pet ID is required")
	}
	if o.UserID <= 0 {
		return errors.New("user ID is required")
	}
	if o.Quantity != 1 {
		return errors.New("quantity must be 1")
	}
	return nil
}
