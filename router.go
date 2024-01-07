package filedispatch

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type destination struct {
	url     string
	pattern []*regexp.Regexp
}

type config struct {
	path        string
	destination []destination
}

func (c *config) getDestination(path string) string {
	return ""
}

func (c *config) addDestination(dest destination) {
	c.destination = append(c.destination, dest)
}

type ConfigBasedRouter struct {
	cfg config
}

func (c ConfigBasedRouter) getConfig(path string) (config, error) {
	return config{}, nil
}

func (c ConfigBasedRouter) Route(path string) error {
	return nil
}

func (c ConfigBasedRouter) WithConfigFile(path string) (ConfigBasedRouter, error) {
	var err error
	res, err := c.getConfig(path)

	if err != nil {
		return ConfigBasedRouter{}, err
	}

	c.cfg = res
	return c, nil
}

func NewRouter() (ConfigBasedRouter, error) {
	var err error

	path, err := getConfigPath()

	if err != nil {
		return ConfigBasedRouter{}, err
	}

	r := ConfigBasedRouter{}

	r, err = r.WithConfigFile(path)

	if err != nil {
		return ConfigBasedRouter{}, err
	}

	return r, nil
}

func getConfigPath() (string, error) {
	var err error
	filename := ".filedispatch/config"
	isXDG := true

	dir, err := os.UserConfigDir()

	if err != nil {
		isXDG = false
		dir, err = os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("cannot get the configuration directory (reason: %v)", err)
		}
	}

	if isXDG {
		filename = strings.TrimPrefix(filename, ".")
	}

	return filepath.Join(dir, filename), nil
}
