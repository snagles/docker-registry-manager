package manager

import (
	"fmt"

	"github.com/heroku/docker-registry-client/registry"
)

// Registries contains a map of all active registries identified by their name
var Registries map[string]*registry.Registry

func init() {
	Registries = map[string]*registry.Registry{}
}

func AddRegistry(scheme, host string, port int) error {
	url := fmt.Sprintf(fmt.Sprintf("%s://%s:%v", scheme, host, port))
	hub, err := registry.New(url, "", "")
	if err != nil {
		return err
	}
	Registries[hub.URL] = hub
	return nil
}
