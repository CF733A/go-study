package pet

import (
	"context"
	"encoding/json"
	"fmt"
	"pet-store/internal/entities"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PetRepository struct {
	db *pgxpool.Pool
}

func NewPetRepository(db *pgxpool.Pool) *PetRepository {
	return &PetRepository{db: db}
}

func (pr *PetRepository) Create(ctx context.Context, pet *entities.Pet) error {
	query := `
		INSERT INTO pets (name, category_id, category_name, photo_urls, tags, status) 
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	
	tagsJSON, err := json.Marshal(pet.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}
	
	err = pr.db.QueryRow(ctx, query,
		pet.Name,
		pet.Category.ID,
		pet.Category.Name,
		pet.PhotoURLs,
		tagsJSON,
		pet.Status,
	).Scan(&pet.ID)

	if err != nil {
		return fmt.Errorf("failed to create pet: %w", err)
	}

	return nil
}

func (pr *PetRepository) GetByID(ctx context.Context, id int64) (*entities.Pet, error) {
	query := `
		SELECT id, name, category_id, category_name, photo_urls, tags, status
		FROM pets
		WHERE id = $1
	`

	pet := &entities.Pet{}
	var tagsJSON []byte
	
	err := pr.db.QueryRow(ctx, query, id).Scan(
		&pet.ID,
		&pet.Name,
		&pet.Category.ID,
		&pet.Category.Name,
		&pet.PhotoURLs,
		&tagsJSON,
		&pet.Status,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get pet by ID: %w", err)
	}

	err = json.Unmarshal(tagsJSON, &pet.Tags)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
	}

	return pet, nil
}

func (pr *PetRepository) Update(ctx context.Context, pet *entities.Pet) error {
	query := `
		UPDATE pets
		SET name = $1, category_id = $2, category_name = $3, photo_urls = $4, tags = $5, status = $6
		WHERE id = $7
	`
	
	tagsJSON, err := json.Marshal(pet.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}
	
	_, err = pr.db.Exec(ctx, query,
		pet.Name,
		pet.Category.ID,
		pet.Category.Name,
		pet.PhotoURLs,
		tagsJSON,
		pet.Status,
		pet.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update pet: %w", err)
	}

	return nil
}

func (pr *PetRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM pets WHERE id = $1`

	_, err := pr.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete pet: %w", err)
	}

	return nil
}

func (pr *PetRepository) UpdateWithFormData(ctx context.Context, id int64, name string, status string) error {
	query := `
		UPDATE pets
		SET name = $1, status = $2
		WHERE id = $3
	`

	_, err := pr.db.Exec(ctx, query, name, status, id)
	if err != nil {
		return fmt.Errorf("failed to update pet with form data: %w", err)
	}

	return nil
}

func (pr *PetRepository) FindByStatus(ctx context.Context, status entities.PetStatus) ([]*entities.Pet, error) {
	query := `
		SELECT id, name, category_id, category_name, photo_urls, tags, status
		FROM pets
		WHERE status = $1
		ORDER BY name
	`
	rows, err := pr.db.Query(ctx, query, status)
	if err != nil {
		return nil, fmt.Errorf("failed to find pets by status: %w", err)
	}
	defer rows.Close()

	var pets []*entities.Pet
	for rows.Next() {
		pet := &entities.Pet{}
		var tagsJSON []byte
		
		err := rows.Scan(
			&pet.ID,
			&pet.Name,
			&pet.Category.ID,
			&pet.Category.Name,
			&pet.PhotoURLs,
			&tagsJSON,
			&pet.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan pet: %w", err)
		}

		err = json.Unmarshal(tagsJSON, &pet.Tags)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
		}
		
		pets = append(pets, pet)
	}
	return pets, nil
}

func (pr *PetRepository) AddImage(ctx context.Context, id int64, imageURL string) error {
	query := `
		UPDATE pets 
		SET photo_urls = array_append(photo_urls, $1)
		WHERE id = $2
	`

	_, err := pr.db.Exec(ctx, query, imageURL, id)
	if err != nil {
		return fmt.Errorf("failed to add image to pet: %w", err)
	}

	return nil
}
