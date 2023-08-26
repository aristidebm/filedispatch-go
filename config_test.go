package main

import (
	"testing"
)

func TestConfigDestinationCanBeEmpty(t *testing.T) {

	config := `
source: /tmp
destination: []`
	configObj, err := ParseConfig(config)

	if err != nil {
		t.Errorf("Got the error %q when parsing the config file", err)
	}

	expectedSource := "/tmp"

	if configObj.Source != expectedSource {
		t.Errorf("Expected config.source=%q, but got config.source=%q", expectedSource, configObj.Source)
	}

	expectedDestination := []Destination{}

	if cap(configObj.Destination) != cap(expectedDestination) {
		t.Errorf("Expected config.destination=%q, but got config.destination=%q", expectedDestination, configObj.Destination)
	}
}

func TestParseConfigWithNonEmptyDestination(t *testing.T) {

	config := `
source: /tmp
destination:
  - path: /home/user/Videos
    pattern: ["*.mp4", "*.webm"]`

	configObj, err := ParseConfig(config)

	if err != nil {
		t.Errorf("Got the error %q when parsing the config file", err)
	}

	expectedSource := "/tmp"

	if configObj.Source != expectedSource {
		t.Errorf("Expected config.source=%q, but got config.source=%q", expectedSource, configObj.Source)
	}

	d := configObj.Destination[0]

	expectedPath := "/home/user/Videos"

	if d.Path != expectedPath {
		t.Errorf("Expected the Path of %v to be %s, but got %s", d, expectedPath, d.Path)
	}

	if len(configObj.Destination) == 0 {
		t.Error("Expected config destination to be not empty")
	}

	if d.Path != expectedPath {
		t.Errorf("Expected config destination %v path to be %s, but got %s", d, expectedPath, d.Path)
	}

	if len(d.Pattern) == 0 {
		t.Errorf("Expected destination %v pattern to be not empty", d)
	}
}

func TestConfigCanContainDestinationWithEmptyPattern(t *testing.T) {

	config := `
source: /tmp
destination:
  - path: /home/user/Videos
    pattern: []`

	configObj, err := ParseConfig(config)

	if err != nil {
		t.Errorf("Got the error %q when parsing the config file", err)
	}

	expectedSource := "/tmp"

	if configObj.Source != expectedSource {
		t.Errorf("Expected config.source=%q, but got config.source=%q", expectedSource, configObj.Source)
	}

	d := configObj.Destination[0]

	expectedPath := "/home/user/Videos"

	if d.Path != expectedPath {
		t.Errorf("Expected the Path of %v to be %s, but got %s", d, expectedPath, d.Path)
	}

	if len(configObj.Destination) == 0 {
		t.Error("Expected config destination to be not empty")
	}

	if d.Path != expectedPath {
		t.Errorf("Expected config destination %v path to be %s, but got %s", d, expectedPath, d.Path)
	}

	expectedPattern := []string{}

	if len(d.Pattern) != len(expectedPattern) {
		t.Errorf("Expected destination %v pattern to be empty", d)
	}
}

func TestConfigMustContainASource(t *testing.T) {

	config := `
destination:
  - path: /home/user/Videos
    pattern: []`

	_, err := ParseConfig(config)

	if err == nil {
		t.Errorf("Expected an error")
	}

}

func TestConfigDestinationMustContainAPath(t *testing.T) {

	config := `
source: /tmp
destination:
  - path: /home/user/Videos
    pattern: [*.mp4]`

	_, err := ParseConfig(config)

	if err == nil {
		t.Errorf("Expected an error")
	}

}

func TestConfigSourceDirPathMustExists(t *testing.T) {

	config := `
source: /unknown
destination:
  - path: /home/user/Videos
    pattern: [*.mp4]`

	_, err := ParseConfig(config)

	if err == nil {
		t.Errorf("Expected an error")
	}
}
