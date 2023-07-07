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
		err = errors.New("request was not authenticated")
		log.Println(err)

		forbidden(w)
		return
	}
	var task models.TaskRequest
	err = task.Process(t.config, r)

	// TODO: Add error translation layer
	if err != nil {
		// FIXME: Get this code from an error struct
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})

		return
	}

	w.WriteHeader(http.StatusOK)

	// TODO: Use a struct for the response structures
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}
