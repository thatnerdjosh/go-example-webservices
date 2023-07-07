package models

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
)

type TaskRequest struct {
	Name string `json:"name"`
}

func (t *TaskRequest) Process(cfg *config.TaskConfig, r *http.Request) error {
	b, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	err := t.load(b)
	if err != nil {
		return err
	}

	return t.run(cfg)
}

func (t *TaskRequest) load(input []byte) error {
	err := json.Unmarshal(input, &t)
	if t.Name == "" {
		err = fmt.Errorf("malformed request, received %v", t)
		log.Println(err)
	}

	return err
}

func (t TaskRequest) run(cfg *config.TaskConfig) error {
	taskItem := cfg.GetTask(t.Name)
	if taskItem == nil {
		return errors.New("")
		// Task not found. Return 404
		// w.WriteHeader(http.StatusNotFound)
		// json.NewEncoder(w).Encode(map[string]interface{}{
		// 	"error":   "Task not found",
		// 	"success": false,
		// })
		// return
	}

	// NOTE: typically not recommended, but should be ok since it is driven entirely by a config file (no user input).
	cmdParts := strings.Split(taskItem.Command, " ")
	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)

	_, err := cmd.Output()
	if err != nil {
		fmt.Println("could not run command: ", err)
		return err
	}

	return nil
}
