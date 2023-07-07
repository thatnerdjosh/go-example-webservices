package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

const defaultConfigPath = "config/tasks.yaml"

type TaskConfig struct {
	Items []TaskConfigItem `yaml:"tasks"`
	Path  string
}

type TaskConfigItem struct {
	Name    string `yaml:"name"`
	Command string `yaml:"command"`
}

func (t *TaskConfig) MustLoad() {
	if t.Path == "" {
		t.Path = defaultConfigPath
	}

	yamlFile, err := ioutil.ReadFile(t.Path)
	if err != nil {
		log.Fatalf("error: %v ", err)
	}

	err = yaml.Unmarshal(yamlFile, &t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

// GetTask returns the task with the given name or nil
func (t TaskConfig) GetTask(name string) *TaskConfigItem {
	for _, task := range t.Items {
		if task.Name == name {
			return &task
		}
	}

	return nil
}
