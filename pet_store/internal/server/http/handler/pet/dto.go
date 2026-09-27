package pet

import "pet-store/internal/entities"

// UploadImageRequest represents the request structure for uploading pet images
// @Description Request payload for uploading a pet image URL
type UploadImageRequest struct {
	ImageURL string `json:"imageUrl" example:"https://example.com/pet-image.jpg"`
}

// CreatePetRequest represents the request structure for creating a new pet
// @Description Request payload for creating a new pet in the store
type CreatePetRequest struct {
	Name      string             `json:"name" example:"doggie"`
	Category  entities.Category  `json:"category"`
	PhotoURLs []string           `json:"photoUrls" example:"[\"https://example.com/photo1.jpg\",\"https://example.com/photo2.jpg\"]"`
	Tags      []entities.Tag     `json:"tags"`
	Status    entities.PetStatus `json:"status" example:"available"`
}

// UpdatePetRequest represents the request structure for updating an existing pet
// @Description Request payload for updating pet details
type UpdatePetRequest struct {
	ID        int64              `json:"id" example:"1"`
	Name      string             `json:"name" example:"doggie"`
	Category  entities.Category  `json:"category"`
	PhotoURLs []string           `json:"photoUrls" example:"[\"https://example.com/photo1.jpg\",\"https://example.com/photo2.jpg\"]"`
	Tags      []entities.Tag     `json:"tags"`
	Status    entities.PetStatus `json:"status" example:"available"`
}