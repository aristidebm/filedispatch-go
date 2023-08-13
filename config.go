package main

import (
	"encoding/json"

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

func ParseConfig(config string) (*Config, error) {

	configObj := Config{}

	configByte := []byte(config)

	if err := json.Unmarshal(configByte, &configObj); err != nil {
		return &configObj, err
	}

	validate := validator.New()

	if err := validate.Struct(&configObj); err != nil {
		return &configObj, err
	}

	return &configObj, nil
}
