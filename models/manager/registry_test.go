package manager

import (
	"testing"
	"time"
)

func TestAddRegistry(t *testing.T) {
	registry, err := AddRegistry("http", "localhost", 5000, 10*time.Second)
	if registry.URL != "http://localhost:5000" {
		t.Error("Incorrect url for created registry")
	}
	if registry.TTL != (10 * time.Second) {
		t.Error("Incorrect TTL for registry")
	}
	if err != nil || AllRegistries.Registries["http://localhost:5000"] != registry {
		t.Error("Unsuccessful in adding a new registry")
	}
}
