package pet

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"pet-store/internal/entities"
	"pet-store/internal/server/http/response"
	services "pet-store/internal/services"

	"github.com/go-chi/chi/v5"
)

type PetHandler struct {
	petService services.PetService
}

func NewPetHandler(petService services.PetService) *PetHandler {
	return &PetHandler{
		petService: petService,
	}
}


// AddPet godoc
// @Summary Add a new pet to the store
// @Description Add a new pet with the provided details
// @Tags pet
// @Accept json
// @Produce json
// @Param pet body CreatePetRequest true "Pet object that needs to be added"
// @Success 201 {object} entities.Pet "Pet created successfully"
// @Failure 400 {object} response.Response "Invalid input or missing required fields"
// @Failure 405 {object} response.Response "Method not allowed"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /pet [post]
// @Example request
// {
//   "name": "doggie",
//   "category": {"id": 1, "name": "Dogs"},
//   "photoUrls": ["url1", "url2"],
//   "tags": [{"id": 1, "name": "friendly"}],
//   "status": "available"
// }
func (h *PetHandler) AddPet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.JSONer(w, http.StatusMethodNotAllowed, nil, errors.New("method not allowed"))
		return
	}

	var req CreatePetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSONer(w, http.StatusBadRequest, nil, errors.New("invalid JSON format"))
		return
	}

	if req.Name == "" {
		response.JSONer(w, http.StatusBadRequest, nil, errors.New("pet name is required"))
		return
	}

	pet := &entities.Pet{
		Name:      req.Name,
		Category:  req.Category,
		PhotoURLs: req.PhotoURLs,
		Tags:      req.Tags,
		Status:    req.Status,
	}

	if pet.Status == "" {
		pet.Status = entities.StatusAvailable
	}

	if err := h.petService.CreatePet(r.Context(), pet); err != nil {
		response.JSONer(w, http.StatusInternalServerError, nil, err)
		return
	}

	response.JSONer(w, http.StatusCreated, pet, nil)
}

// UpdatePet godoc
// @Summary Update an existing pet
// @Description Update pet details by ID
// @Tags pet
// @Accept json
// @Produce json
// @Param pet body UpdatePetRequest true "Pet object with updated details"
// @Success 200 {object} entities.Pet "Pet updated successfully"
// @Failure 400 {object} response.Response "Invalid ID supplied or invalid input"
// @Failure 404 {object} response.Response "Pet not found"
// @Failure 405 {object} response.Response "Method not allowed"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /pet [put]
// @Example response
// {
//   "id": 1,
//   "name": "doggie",
//   "category": {"id": 1, "name": "Dogs"},
//   "photoUrls": ["url1", "url2"],
//   "tags": [{"id": 1, "name": "friendly"}],
//   "status": "available"
// }
func (h *PetHandler) UpdatePet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		response.JSONer(w, http.StatusMethodNotAllowed, nil, errors.New("method not allowed"))
		return
	}

	var req UpdatePetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSONer(w, http.StatusBadRequest, nil, errors.New("invalid JSON format"))
		return
	}

	if req.ID <= 0 {
		response.JSONer(w, http.StatusBadRequest, nil, errors.New("pet ID is required"))
		return
	}
	if req.Name == "" {
		response.JSONer(w, http.StatusBadRequest, nil, errors.New("pet name is required"))
		return
	}

	pet := &entities.Pet{
		ID:        req.ID,
		Name:      req.Name,
		Category:  req.Category,
		PhotoURLs: req.PhotoURLs,
		Tags:      req.Tags,
		Status:    req.Status,
	}

	if err := h.petService.UpdatePet(r.Context(), pet); err != nil {
		if err.Error() == "pet not found" {
			response.JSONer(w, http.StatusNotFound, nil, err)
		} else if err.Error() == "cannot update sold pet" {
			response.JSONer(w, http.StatusBadRequest, nil, err)
		} else {
			response.JSONer(w, http.StatusInternalServerError, nil, err)
		}
		return
	}

	response.JSONer(w, http.StatusOK, pet, nil)
}

// FindPetByID godoc
// @Summary Find pet by ID
// @Description Returns a single pet by ID
// @Tags pet
// @Produce json
// @Param petId path int64 true "ID of pet to return"
// @Success 200 {object} entities.Pet "Successful operation"
// @Failure 400 {object} response.Response "Invalid ID supplied"
// @Failure 404 {object} response.Response "Pet not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /pet/{petId} [get]
// @Example response
// {
//   "id": 1,
//   "name": "doggie",
//   "category": {"id": 1, "name": "Dogs"},
//   "photoUrls": ["url1", "url2"],
//   "tags": [{"id": 1, "name": "friendly"}],
//   "status": "available"
// }
func (h *PetHandler) FindPetByID(w http.ResponseWriter, r *http.Request) {
	petIDStr := chi.URLParam(r, "petId")
	if petIDStr == "" {
		response.JSONer(w, http.StatusBadRequest, nil, errors.New("pet ID is required"))
		return
	}

	petID, err := strconv.ParseInt(petIDStr, 10, 64)
	if err != nil || petID <= 0 {
		response.JSONer(w, http.StatusBadRequest, nil, errors.New("invalid pet ID"))
		return
	}

	pet, err := h.petService.GetPetByID(r.Context(), petID)
	if err != nil {
		if err.Error() == "pet not found" {
			response.JSONer(w, http.StatusNotFound, nil, err)
		} else {
			response.JSONer(w, http.StatusInternalServerError, nil, err)
		}
		return
	}

	response.JSONer(w, http.StatusOK, pet, nil)
}

// UpdatePetWithFormData godoc
// @Summary Updates a pet in the store with form data
// @Description Update pet name and/or status using form data
// @Tags pet
// @Accept multipart/form-data
// @Produce json
// @Param petId path int64 true "ID of pet that needs to be updated"
// @Param name formData string false "Updated name of the pet"
// @Param status formData string false "Updated status of the pet"
// @Success 200 {object} response.Response "Pet updated successfully"
// @Failure 400 {object} response.Response "Invalid input"
// @Failure 404 {object} response.Response "Pet not found"
// @Failure 405 {object} response.Response "Method not allowed"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /pet/{petId} [post]
// @Example response
// {
//   "status": true,
//   "message": "pet updated successfully"
// }
func (h *PetHandler) UpdatePetWithFormData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.JSONer(w, http.StatusMethodNotAllowed, nil, errors.New("method not allowed"))
		return
	}

	petIDStr := chi.URLParam(r, "petId")
	if petIDStr == "" {
		response.JSONer(w, http.StatusBadRequest, nil, errors.New("pet ID is required"))
		return
	}

	petID, err := strconv.ParseInt(petIDStr, 10, 64)
	if err != nil || petID <= 0 {
		response.JSONer(w, http.StatusBadRequest, nil, errors.New("invalid pet ID"))
		return
	}

	if err := r.ParseForm(); err != nil {
		response.JSONer(w, http.StatusBadRequest, nil, errors.New("invalid form data"))
		return
	}

	name := r.FormValue("name")
	status := r.FormValue("status")

	if name == "" && status == "" {
		response.JSONer(w, http.StatusBadRequest, nil, errors.New("at least one field (name or status) is required"))
		return
	}

	if err := h.petService.UpdatePetWithForm(r.Context(), petID, name, status); err != nil {
		if err.Error() == "pet not found" {
			response.JSONer(w, http.StatusNotFound, nil, err)
		} else {
			response.JSONer(w, http.StatusInternalServerError, nil, err)
		}
		return
	}

	response.JSONer(w, http.StatusOK, map[string]string{"message": "pet updated successfully"}, nil)
}

// DeletePet godoc
// @Summary Delete a pet
// @Description Deletes a pet by ID
// @Tags pet
// @Produce json
// @Param petId path int64 true "Pet ID to delete"
// @Success 200 {object} response.Response "Pet deleted successfully"
// @Failure 400 {object} response.Response "Invalid ID supplied"
// @Failure 404 {object} response.Response "Pet not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /pet/{petId} [delete]
// @Example response
// {
//   "status": true,
//   "message": "pet deleted successfully"
// }
func (h *PetHandler) DeletePet(w http.ResponseWriter, r *http.Request) {
	petIDStr := chi.URLParam(r, "petId")
	if petIDStr == "" {
		response.JSONer(w, http.StatusBadRequest, nil, errors.New("pet ID is required"))
		return
	}

	petID, err := strconv.ParseInt(petIDStr, 10, 64)
	if err != nil || petID <= 0 {
		response.JSONer(w, http.StatusBadRequest, nil, errors.New("invalid pet ID"))
		return
	}

	if err := h.petService.DeletePet(r.Context(), petID); err != nil {
		if err.Error() == "pet not found" {
			response.JSONer(w, http.StatusNotFound, nil, err)
		} else {
			response.JSONer(w, http.StatusInternalServerError, nil, err)
		}
		return
	}

	response.JSONer(w, http.StatusOK, map[string]string{"message": "pet deleted successfully"}, nil)
}

// UploadImage godoc
// @Summary Uploads an image for a pet
// @Description Add an image URL to a pet's photo gallery
// @Tags pet
// @Accept json
// @Produce json
// @Param petId path int64 true "ID of pet to update"
// @Param image body UploadImageRequest true "Image URL object"
// @Success 200 {object} response.Response "Image uploaded successfully"
// @Failure 400 {object} response.Response "Invalid input"
// @Failure 404 {object} response.Response "Pet not found"
// @Failure 405 {object} response.Response "Method not allowed"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /pet/{petId}/uploadImage [post]
// @Example response
// {
//   "status": true,
//   "message": "image uploaded successfully"
// }
func (h *PetHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.JSONer(w, http.StatusMethodNotAllowed, nil, errors.New("method not allowed"))
		return
	}

	petIDStr := chi.URLParam(r, "petId")
	if petIDStr == "" {
		response.JSONer(w, http.StatusBadRequest, nil, errors.New("pet ID is required"))
		return
	}

	petID, err := strconv.ParseInt(petIDStr, 10, 64)
	if err != nil || petID <= 0 {
		response.JSONer(w, http.StatusBadRequest, nil, errors.New("invalid pet ID"))
		return
	}

	var req UploadImageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSONer(w, http.StatusBadRequest, nil, errors.New("invalid JSON format"))
		return
	}

	if req.ImageURL == "" {
		response.JSONer(w, http.StatusBadRequest, nil, errors.New("image URL is required"))
		return
	}

	if err := h.petService.UploadPetImage(r.Context(), petID, req.ImageURL); err != nil {
		if err.Error() == "pet not found" {
			response.JSONer(w, http.StatusNotFound, nil, err)
		} else {
			response.JSONer(w, http.StatusInternalServerError, nil, err)
		}
		return
	}

	response.JSONer(w, http.StatusOK, map[string]string{"message": "image uploaded successfully"}, nil)
}

// FindPetsByStatus godoc
// @Summary Finds pets by status
// @Description Multiple status values can be provided with comma separated strings
// @Tags pet
// @Produce json
// @Param status query string true "Status values that need to be considered for filter"
// @Success 200 {array} entities.Pet "Successful operation"
// @Failure 400 {object} response.Response "Invalid status value"
// @Failure 405 {object} response.Response "Method not allowed"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /pet/findByStatus [get]
// @Example response
// [
//   {
//     "id": 1,
//     "name": "doggie",
//     "category": {"id": 1, "name": "Dogs"},
//     "photoUrls": ["url1", "url2"],
//     "tags": [{"id": 1, "name": "friendly"}],
//     "status": "available"
//   }
// ]
func (h *PetHandler) FindPetsByStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.JSONer(w, http.StatusMethodNotAllowed, nil, errors.New("method not allowed"))
		return
	}

	statusValues := r.URL.Query()["status"]
	if len(statusValues) == 0 {
		response.JSONer(w, http.StatusBadRequest, nil, errors.New("status parameter is required"))
		return
	}

	status := entities.PetStatus(statusValues[0])

	pets, err := h.petService.FindPetsByStatus(r.Context(), status)
	if err != nil {
		response.JSONer(w, http.StatusInternalServerError, nil, err)
		return
	}

	response.JSONer(w, http.StatusOK, pets, nil)
}