package home

type HomeResponse struct {
    Service       string            `json:"service"`
    Status        string            `json:"status"`
    Message       string            `json:"message"`
    Documentation string            `json:"documentation"`
    Endpoints     map[string]string `json:"endpoints"`
}

func NewHomeResponse() HomeResponse {
    return HomeResponse{
        Service:       "geo-service",
        Status:        "ok",
        Message:       "Добро пожаловать в GEO Service API",
        Documentation: "http://localhost:8080/swagger/",
        Endpoints: map[string]string{
            "register": "/api/register",
            "login":    "/api/login",
            "search":   "/api/address/search (protected)",
            "geocode":  "/api/address/geocode (protected)",
        },
    }
}