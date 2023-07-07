package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type TaskConfig struct {
	Items []TaskConfigItem `yaml:"tasks"`
}

type TaskConfigItem struct {
	Name    string `yaml:"name"`
	Command string `yaml:"command"`
}

func (t *TaskConfig) MustLoad() {
	yamlFile, err := ioutil.ReadFile("config/tasks.yaml")
	if err != nil {
		log.Fatalf("error: %v ", err)
	}

	err = yaml.Unmarshal(yamlFile, &t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
