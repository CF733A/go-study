package home

// HomeResponse represents API home page information
// @Description Service information response structure
type HomeResponse struct {
    Service       string `json:"service" example:"pet-store"`
    Status        string `json:"status" example:"ok"`
    Message       string `json:"message" example:"Добро пожаловать в PetStore"`
    Documentation string `json:"documentation" example:"http://localhost:8080/swagger/"`
}

func NewHomeResponse() HomeResponse {
    return HomeResponse{
        Service:       "pet-store",
        Status:        "ok",
        Message:       "Добро пожаловать в PetStore",
        Documentation: "http://localhost:8080/swagger/",
    }
}