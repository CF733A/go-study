package responder

import (
    "encoding/json"
    "log"
    "net/http"
)

type Responder interface {
    Success(w http.ResponseWriter, data interface{})
    Error(w http.ResponseWriter, status int, message string)
    BadRequest(w http.ResponseWriter, message string)
    Unauthorized(w http.ResponseWriter, message string)
    InternalError(w http.ResponseWriter, message string)
}

type JSONResponder struct{}

type Response struct {
    Status string      `json:"status"`
    Data   interface{} `json:"data,omitempty"`
    Error  string      `json:"error,omitempty"`
}

func NewJSONResponder() Responder {
    return &JSONResponder{}
}

func (j *JSONResponder) Success(w http.ResponseWriter, data interface{}) {
    resp := &Response{
        Status: "success",
        Data:   data,
    }

    j.sendJSON(w, http.StatusOK, resp)
}

func (j *JSONResponder) Error(w http.ResponseWriter, status int, message string) {
    resp := &Response{
        Status: "error",
        Error:  message,
    }

    j.sendJSON(w, status, resp)
}

func (j *JSONResponder) BadRequest(w http.ResponseWriter, message string) {
    j.Error(w, http.StatusBadRequest, message)
}

func (j *JSONResponder) Unauthorized(w http.ResponseWriter, message string) {
    j.Error(w, http.StatusUnauthorized, message)
}

func (j *JSONResponder) InternalError(w http.ResponseWriter, message string) {
    j.Error(w, http.StatusInternalServerError, message)
}

func (j *JSONResponder) sendJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.WriteHeader(status)
    encoder := json.NewEncoder(w)
    encoder.SetIndent("", "  ")
    if err := encoder.Encode(data); err != nil {
        log.Printf("JSON encode error: %v", err)
        return
    }
}