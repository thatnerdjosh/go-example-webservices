package controllers

import (
	"encoding/json"
	"net/http"
)

func badRequest(w http.ResponseWriter) {
	// Data wasn't able to unmarshal, return with 400
	w.WriteHeader(http.StatusBadRequest)

	// TODO: Add error translation layer
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":   "Malformed request",
		"success": false,
	})
}

func forbidden(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":   "Not authenticated",
		"success": false,
	})
}

func authenticated(r *http.Request) bool {
	// Authentication
	// TODO: HTTP request to auth server
	return r.Header.Get("Authorization") != ""
}
