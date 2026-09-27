package pet

import "github.com/go-chi/chi/v5"

func (h *PetHandler) RegisterProtectedRoutes(r chi.Router) {
    r.Route("/pet", func(r chi.Router) {
        r.Post("/", h.AddPet)
        r.Put("/", h.UpdatePet)
        r.Get("/findByStatus", h.FindPetsByStatus)
        r.Route("/{petId}", func(r chi.Router) {
            r.Get("/", h.FindPetByID)
            r.Post("/", h.UpdatePetWithFormData)
            r.Delete("/", h.DeletePet)
            r.Post("/uploadImage", h.UploadImage)
        })
    })
}