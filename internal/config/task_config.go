package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

const defaultConfigPath = "config/tasks.yaml"

type TaskConfig struct {
	Dependencies map[string]Dependency
	Tasks        []Task `yaml:"tasks"`
	ConfigPath   string
}

type Dependency struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	TLS  bool   `yaml:"tls"`
}

type Task struct {
	Name    string `yaml:"name"`
	Command string `yaml:"command"`
}

func (t *TaskConfig) MustLoad() {
	if t.ConfigPath == "" {
		t.ConfigPath = defaultConfigPath
	}

	yamlFile, err := os.Open(t.ConfigPath)
	if err != nil {
		log.Fatalf("error: %v ", err)
	}

	err = yaml.NewDecoder(yamlFile).Decode(&t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	authAPI, ok := t.Dependencies["authAPI"]
	if !ok {
		log.Fatalf("error: authAPI is not configured. Check dependency config.")
	}

	if host, has := os.LookupEnv("AUTH_API_HOST"); has {
		authAPI.Host = host
		fmt.Println(host)
	}

	if portStr, has := os.LookupEnv("AUTH_API_PORT"); has {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			log.Fatalf("error: AUTH_API_PORT was set but was not an integer")
		}

		authAPI.Port = port
	}

	t.Dependencies["authAPI"] = authAPI
}

// GetTask returns the task with the given name or nil
func (t TaskConfig) GetTask(name string) *Task {
	for _, task := range t.Tasks {
		if task.Name == name {
			return &task
		}
	}

	return nil
}

func (t TaskConfig) GetAuthAPIURL() string {
	authAPI := t.Dependencies["authAPI"]
	proto := "https"
	if !authAPI.TLS {
		proto = "http"
	}

	return fmt.Sprintf("%s://%s:%d", proto, authAPI.Host, authAPI.Port)
}
