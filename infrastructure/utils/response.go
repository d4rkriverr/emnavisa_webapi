package utils

import (
	"encoding/json"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, statusCode int, message string) {
	response := map[string]any{
		"success": false,
		"message": message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// RespondWithSuccess writes a JSON success response to the client
func RespondWithSuccess(w http.ResponseWriter, data interface{}) {
	response := map[string]any{
		"success": true,
		"payload": data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
