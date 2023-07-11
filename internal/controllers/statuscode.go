package controllers

import (
	"encoding/json"
	"net/http"
)

func forbidden(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":   "Not authenticated",
		"success": false,
	})
}
