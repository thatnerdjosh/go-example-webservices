package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"

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

	b, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	// if err != nil {
	// 	// Error processing request (e.g., timeout), return 500
	// 	http.Error(w, err.Error(), 500)
	// 	return
	// }

	// TODO: Extract logic to models
	var task models.TaskRequest
	err = json.Unmarshal(b, &task)
	if task.Name == "" {
		err = fmt.Errorf("malformed request, received %v", task)
		log.Println(err)
	}

	if err != nil {
		badRequest(w)
		return
	}

	taskItem := t.config.GetTask(task.Name)
	if taskItem == nil {
		// Task not found. Return 404
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   "Task not found",
			"success": false,
		})
		return
	}

	// NOTE: typically not recommended, but should be ok since it is driven entirely by a config file (no user input).
	cmdParts := strings.Split(taskItem.Command, " ")
	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)

	_, err = cmd.Output()
	if err != nil {
		fmt.Println("could not run command: ", err)
	}

	w.WriteHeader(http.StatusOK)

	// TODO: Use a struct for the response structures
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}

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
