package manager

import (
	"fmt"
	"sync"
	"time"

	client "github.com/heroku/docker-registry-client/registry"
)

var AllRegistries Registries

// Registries contains a map of all active registries identified by their name
type Registries struct {
	Registries map[string]*Registry
	sync.Mutex
}

type Registry struct {
	*client.Registry
	Repositories map[string]*Repository
	TTL          time.Duration
	sync.Mutex
}

func (r *Registry) Refresh() {

	refreshedRegistry := Registry{}
	repoList, _ := refreshedRegistry.Registry.Repositories()
	// Get the repository information
	refreshedRegistry.Repositories = map[string]*Repository{}
	for _, repoName := range repoList {

		// Build a repository object
		repo := Repository{Name: repoName}

		tagList, _ := refreshedRegistry.Tags(repoName)
		repo.Tags = map[string]*Tag{}
		for _, tagName := range tagList {
			man, _ := refreshedRegistry.Manifest(repoName, tagName)
			tag := Tag{
				Name:           tagName,
				SignedManifest: man,
			}

			// Add the tag to the repository
			repo.Tags[tagName] = &tag
		}
		r.Repositories[repoName] = &repo
	}
}

func (r *Registry) Status() string {
	if err := r.Ping(); err != nil {
		return "DOWN"
	}
	return "UP"
}

func init() {
	AllRegistries.Registries = map[string]*Registry{}
}

func AddRegistry(scheme, host string, port int, ttl time.Duration) (*Registry, error) {
	url := fmt.Sprintf(fmt.Sprintf("%s://%s:%v", scheme, host, port))
	hub, err := client.New(url, "", "")
	if err != nil {
		return nil, err
	}
	r := Registry{
		Registry: hub,
		TTL:      ttl,
	}

	AllRegistries.Lock()
	AllRegistries.Registries[hub.URL] = &r
	AllRegistries.Unlock()
	return &r, nil
}
