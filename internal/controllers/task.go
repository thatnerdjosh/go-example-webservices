package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/thatnerdjosh/example-webservices/internal/config"
	"github.com/thatnerdjosh/example-webservices/internal/models"
)

type TaskController struct {
	config *config.TaskConfig
}

func NewTaskController(taskConfig *config.TaskConfig) TaskController {
	var ctrl TaskController

	ctrl.config = &config.TaskConfig{}
	if taskConfig != nil {
		ctrl.config = taskConfig
	}

	ctrl.config.MustLoad()
	return ctrl
}

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

	// if err != nil {
	// 	// Error processing request (e.g., timeout), return 500
	// 	http.Error(w, err.Error(), 500)
	// 	return
	// }

	var task models.TaskRequest
	err = task.Process(t.config, r)

	// TODO: Extract to separate handler to reduce complexity of this one.
	if err != nil {
		switch {
		case errors.Is(err, models.ErrTaskNotFound):
			// Task not found. Return 404
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error":   err.Error(),
				"success": false,
			})

			return
		case errors.Is(err, models.ErrBadData):
			// Invalid data provided, return 400
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error":   err.Error(),
				"success": false,
			})

			return
		default:
			// Unhandled error, assume ISE for now
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error":   "An unexpected error has occurred.",
				"success": false,
			})

			return
		}
	}

	w.WriteHeader(http.StatusOK)

	// TODO: Use a struct for the response structures
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}
