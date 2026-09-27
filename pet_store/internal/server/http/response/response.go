package response

import (
	"encoding/json"
	"net/http"
)

// Response represents unified JSON response structure
// @Description Unified response structure for all API responses
type Response struct {
	Status  bool        `json:"status" example:"true"`
	Message string      `json:"message,omitempty" example:"Success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty" example:"Error message"`
}

func JSONer(w http.ResponseWriter, statusCode int, data interface{}, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := Response{
		Status: err == nil,
		Data:   data,
	}

	if err != nil {
		response.Error = err.Error()
	} else if data == nil {
		response.Message = "Success"
	}

	json.NewEncoder(w).Encode(response)
}
