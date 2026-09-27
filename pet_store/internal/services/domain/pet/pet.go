package pet

import (
	"context"
	"errors"
	"fmt"
	"pet-store/internal/entities"
	"pet-store/internal/storage"
)

type PetService struct {
	PetRepo storage.PetRepository
}

func NewPetService(petRepo storage.PetRepository) *PetService {
	return &PetService{
		PetRepo: petRepo,
	}
}
func (p *PetService) CreatePet(ctx context.Context, pet *entities.Pet) error {
	if pet.Name == "" {
		return errors.New("pet name required")
	}

	if pet.Status == "" {
		pet.Status = entities.StatusAvailable
	}

	return p.PetRepo.Create(ctx, pet)
}

func (p *PetService) GetPetByID(ctx context.Context, id int64) (*entities.Pet, error) {
	if id <= 0 {
		return nil, errors.New("invalid pet ID")
	}

	pet, err := p.PetRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("pet not found: %w", err)
	}

	return pet, nil
}

func (p *PetService) UpdatePet(ctx context.Context, pet *entities.Pet) error {
	petExists, err := p.PetRepo.GetByID(ctx, pet.ID)
	if err != nil {
		return fmt.Errorf("pet not found: %w", err)
	}

	if petExists.Status == entities.StatusSold {
		return errors.New("cannot update sold pet")
	}

	return p.PetRepo.Update(ctx, pet)
}

func (p *PetService) DeletePet(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("invalid pet ID")
	}

	_, err := p.PetRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("pet not found: %w", err)
	}

	return p.PetRepo.Delete(ctx, id)
}

func (s *PetService) UpdatePetWithForm(ctx context.Context, id int64, name, status string) error {
	if id <= 0 {
		return errors.New("invalid pet ID")
	}

	_, err := s.PetRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("pet not found: %w", err)
	}

	return s.PetRepo.UpdateWithFormData(ctx, id, name, status)
}

func (p *PetService) FindPetsByStatus(ctx context.Context, status entities.PetStatus) ([]*entities.Pet, error) {
	switch status {
	case entities.StatusAvailable, entities.StatusPending, entities.StatusSold:
		return p.PetRepo.FindByStatus(ctx, status)
	default:
		return nil, errors.New("invalid status")
	}
}

func (p *PetService) UploadPetImage(ctx context.Context, id int64, imageURL string) error {
	if id <= 0 {
		return errors.New("invalid pet ID")
	}

	if imageURL == "" {
		return errors.New("image URL is required")
	}

	_, err := p.PetRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("pet not found: %w", err)
	}

	return p.PetRepo.AddImage(ctx, id, imageURL)
}
