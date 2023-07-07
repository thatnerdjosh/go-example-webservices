package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/thatnerdjosh/example-webservices/internal/models"
)

type TaskController struct{}

func (t TaskController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch {
	case r.Method == "POST" && r.URL.Path == "/":
		t.ExecuteTask(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (t TaskController) ExecuteTask(w http.ResponseWriter, r *http.Request) {
	var err error
	if !authenticated(r) {
		err = errors.New("request was not authorized")
		log.Println(err)

		forbidden(w)
		return
	}

	b, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	// if err != nil {
	// 	// Error processing request (e.g., timeout), return 500
	// 	http.Error(w, err.Error(), 500)
	// 	return
	// }

	var task models.Task
	err = json.Unmarshal(b, &task)
	if task.Id == "" {
		err = fmt.Errorf("malformed request, received %v", task)
		log.Println(err)
	}

	if err != nil {
		badRequest(w)
		return
	}

	// if resp == nil {
	// 	// Task not found. Return 404
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }

	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(&resp)
}

func badRequest(w http.ResponseWriter) {
	// Data wasn't able to unmarshal, return with 400
	w.WriteHeader(http.StatusBadRequest)

	// TODO: Add error translation layer
	json.NewEncoder(w).Encode(map[string]string{
		"error": "Malformed request",
	})
}

func forbidden(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(map[string]string{
		"error": "Not authenticated",
	})
}

func authenticated(r *http.Request) bool {
	// Authentication
	// TODO: HTTP request to auth server
	return r.Header.Get("Authorization") != ""
}
