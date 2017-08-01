package manager

import (
	"os"
	"strconv"
	"testing"
)

func TestTagCount(t *testing.T) {
	registry := createTestRegistry()
	if registry.TagCount() != 1 {
		t.Error("Expected 1 tag, received: " + strconv.Itoa(registry.TagCount()))
	}
}

func TestIP(t *testing.T) {
	registry := createTestRegistry()
	if registry.IP() == "" {
		t.Error("Failed to get ip for hostname")
	}
}

func createTestRegistry() Registry {
	hostname, _ := os.Hostname()
	registry := Registry{
		Name:    "test",
		Host:    hostname,
		Scheme:  "https",
		Version: "v2",
		Port:    5000,
	}
	repository := Repository{
		Name: "testrepo",
	}

	repository.Tags = make(map[string]*Tag)
	tag := Tag{
		ID:   "1",
		Name: "testTag",
		Size: 400,
	}
	repository.Tags["testTag"] = &tag
	registry.Repositories = make(map[string]*Repository)
	registry.Repositories["testRepo"] = &repository

	return registry
}
