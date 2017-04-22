package manager

import (
	"testing"
	"time"
)

func TestAddRegistry(t *testing.T) {
	registry, err := AddRegistry("http", "localhost", 5010, 10*time.Second)
	if registry.URL != "http://localhost:5010" {
		t.Error("Incorrect url for created registry")
	}
	if registry.TTL != (10 * time.Second) {
		t.Error("Incorrect TTL for registry")
	}
	if err != nil || AllRegistries.Registries["http://localhost:5010"] != registry {
		t.Error("Unsuccessful in adding a new registry")
	}
}
