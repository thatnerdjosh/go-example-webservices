package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/thatnerdjosh/example-webservices/internal/config"
	"github.com/thatnerdjosh/example-webservices/internal/models"
)

type HttpClient interface {
	Get(string) (*http.Response, error)
}

type TaskResponse struct {
	status  int
	ErrData string `json:"error,omitempty"`
	Success bool   `json:"success"`
}

type TaskController struct {
	config *config.TaskConfig
	client HttpClient
}

var (
	ErrForbidden = errors.New("request was not authenticated")
)

func NewTaskController(taskConfig *config.TaskConfig, client HttpClient) TaskController {
	var ctrl TaskController

	ctrl.config = &config.TaskConfig{}
	if taskConfig != nil {
		ctrl.config = taskConfig
	}

	ctrl.client = client

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
	authenticated, err := t.authenticated(r)
	if err != nil || !authenticated {
		handle(ErrForbidden, w)
		return
	}

	var task models.TaskRequest
	err = task.Process(t.config, r)
	if err != nil {
		handle(err, w)
		return
	}

	TaskResponse{
		status:  http.StatusOK,
		Success: true,
	}.WriteJSON(w)
}

func (t TaskController) authenticated(r *http.Request) (bool, error) {
	if r.Header.Get("Authorization") == "" {
		return false, nil
	}

	resp, err := t.client.Get(t.config.GetAuthAPIURL())
	if err != nil {
		return false, err
	}

	return resp.StatusCode == http.StatusOK, nil
}

func (tr TaskResponse) WriteJSON(w http.ResponseWriter) {
	respData, err := json.Marshal(tr)
	if err != nil {
		log.Println("unable to marshal response data, falling back to ISE.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(tr.status)
	_, err = w.Write(respData)
	if err != nil {
		log.Println(err)
		return
	}
}

func handle(err error, w http.ResponseWriter) {
	log.Println(err)

	// NOTE: We set the ErrData of the struct manually to map to the sentinel error for now
	// this is done to prevent internal error leakage to the outside
	// (consider error wrapping).
	var resp TaskResponse

	switch {
	case errors.Is(err, ErrForbidden):
		// Unauthenticated. Return 403
		resp.status = http.StatusForbidden
		resp.ErrData = ErrForbidden.Error()
		resp.WriteJSON(w)
		return
	case errors.Is(err, models.ErrTaskNotFound):
		// Task not found. Return 404
		resp.status = http.StatusNotFound
		resp.ErrData = models.ErrTaskNotFound.Error()
		resp.WriteJSON(w)
		return
	case errors.Is(err, models.ErrBadData):
		// Invalid data provided, return 400
		resp.status = http.StatusBadRequest
		resp.ErrData = models.ErrBadData.Error()
		resp.WriteJSON(w)
		return
	default:
		// Unhandled error, assume ISE for now
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
