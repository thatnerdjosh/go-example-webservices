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

var (
	ErrTaskNotFound = errors.New("task was not found")
	ErrBadData      = errors.New("malformed data provided")
)

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
	if err != nil {
		log.Println(err)
	}

	if err != nil || t.Name == "" {
		// Error to bubble up to end user
		// TODO: Consider error wrapping along with error handling improvements.
		err = ErrBadData
	}

	return err
}

func (t TaskRequest) run(cfg *config.TaskConfig) error {
	taskItem := cfg.GetTask(t.Name)
	if taskItem == nil {
		return ErrTaskNotFound
	}

	// NOTE: typically not recommended, but should be ok since it is driven entirely by a config file (no user input).
	cmdParts := strings.Split(taskItem.Command, " ")
	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)

	res, err := cmd.Output()
	fmt.Println(string(res))
	if err != nil {
		log.Println("could not run command: ", err)
		return err
	}

	return nil
}
