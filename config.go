package main

import (
	"github.com/ghodss/yaml"
	"github.com/go-playground/validator/v10"
)

type Destination struct {
	Path    string   `json:"path" validate:"required"`
	Pattern []string `json:"pattern"`
}

type Config struct {
	Source      string        `json:"source" validate:"required,dir"`
	Destination []Destination `json:"destination" validate:"required,dive"`
}

func ParseConfig(config string) (Config, error) {
	configObj := Config{}

	if err := yaml.Unmarshal([]byte(config), &configObj); err != nil {
		return configObj, err
	}

	validate := validator.New()

	if err := validate.Struct(&configObj); err != nil {
		return configObj, err
	}

	return configObj, nil
}
