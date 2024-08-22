package utils

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

var currentId int

func MethodResponse(w http.ResponseWriter, r *http.Request, method string) {
	if r.Method != method {
		w.Header().Add("Allow", method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func GetID(w http.ResponseWriter, r *http.Request) int {
	id := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		WriteJsonError(w, "Invalid ID", http.StatusBadRequest, err)
		return 0
	}
	return idInt
}

type Response struct {
	Message string      `json:"message"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Success bool        `json:"success,omitempty"`
}

func WriteJsonResponse(w http.ResponseWriter, data interface{}, message string, status int) {
	response := Response{
		Message: message,
		Success: true,
		Data:    data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

func WriteJsonError(w http.ResponseWriter, message string, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	response := Response{
		Message: message,
		Error:   err.Error(),
		Success: false,
	}
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

func FormatDate(t time.Time) string {
	return t.Format("January 02, 2006")
}

func ParseDate(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}

func GenerateNewID() int {
	currentId++
	return currentId
}
