package order

import (
	"context"
	"errors"
	"fmt"
	"pet-store/internal/entities"
	"pet-store/internal/storage"
)

type OrderService struct {
	OrderRepo storage.OrderRepository
	PetRepo   storage.PetRepository
}

func NewOrderService(orderRepo storage.OrderRepository, petRepo storage.PetRepository) *OrderService {
	return &OrderService{
		OrderRepo: orderRepo,
		PetRepo:   petRepo,
	}
}

func (o *OrderService) CreateOrder(ctx context.Context, order *entities.Order) error {
	if err := order.Validate(); err != nil {
		return fmt.Errorf("order validation failed: %w", err)
	}

	pet, err := o.PetRepo.GetByID(ctx, order.ID)
	if err != nil {
		return errors.New("pet not found")
	}

	if pet.Status != entities.StatusAvailable {
		return errors.New("pet is not available for purchase")
	}

	pet.MarkAsSold()
	if err := o.PetRepo.Update(ctx, pet); err != nil {
		return fmt.Errorf("failed to update pet status: %w", err)
	}

	if err := o.OrderRepo.Create(ctx, order); err != nil {
		pet.Status = entities.StatusAvailable
		_ = o.PetRepo.Update(ctx, pet)
		return fmt.Errorf("failed to create order: %w", err)
	}

	return nil
}

func (o *OrderService) GetOrderByID(ctx context.Context, id int64) (*entities.Order, error) {
	if id <= 0 {
		return nil, errors.New("invalid order ID")
	}

	order, err := o.OrderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	return order, nil
}

func (o *OrderService) DeleteOrder(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("invalid order ID")
	}

	order, err := o.OrderRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	if order.Status == entities.StatusDelivered {
		return errors.New("cannot delete delivered order")
	}

	return o.OrderRepo.Delete(ctx, id)
}

func (o *OrderService) GetInventory(ctx context.Context) (map[string]int32, error) {
	availablePets, err := o.PetRepo.FindByStatus(ctx, entities.StatusAvailable)
	if err != nil {
		return nil, fmt.Errorf("failed to get available pets: %w", err)
	}

	pendingPets, err := o.PetRepo.FindByStatus(ctx, entities.StatusPending)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending pets: %w", err)
	}

	soldPets, err := o.PetRepo.FindByStatus(ctx, entities.StatusSold)
	if err != nil {
		return nil, fmt.Errorf("failed to get sold pets: %w", err)
	}

	inventory := map[string]int32{
		"available": int32(len(availablePets)),
		"pending":   int32(len(pendingPets)),
		"sold":      int32(len(soldPets)),
	}

	return inventory, nil
}