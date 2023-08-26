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

func ParseConfig(config string) (*Config, error) {
	// TODO:
	// https://medium.com/thedevproject/easy-de-serialisation-of-yaml-files-in-go-4557456b0a98
	// + Accept string that contains yaml tags using gopkg.in/yaml/v3.
	// + Convert this string into it's json equivalent.
	// + Pass it to go-playground/validator.

	configObj := Config{}

	// configByte, _ := os.ReadFile("/tmp/config.yml")
	// vf := []byte(config)

	// log.Println(len(vf))
	// log.Println("--------------------------------------------------")
	// log.Println(len(configByte))

	if err := yaml.Unmarshal([]byte(config), &configObj); err != nil {
		return &configObj, err
	}

	validate := validator.New()

	if err := validate.Struct(&configObj); err != nil {
		return &configObj, err
	}

	return &configObj, nil
}
