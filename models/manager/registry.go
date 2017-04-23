package manager

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"code.cloudfoundry.org/bytefmt"

	"github.com/Sirupsen/logrus"
	client "github.com/heroku/docker-registry-client/registry"
	"github.com/snagles/docker-registry-manager/utilities"
)

var AllRegistries Registries

func init() {
	AllRegistries.Registries = map[string]*Registry{}
}

// Registries contains a map of all active registries identified by their name, locked when necessary
type Registries struct {
	Registries map[string]*Registry
	sync.Mutex
}

type Registry struct {
	*client.Registry
	Repositories map[string]*Repository
	TTL          time.Duration
	Ticker       *time.Ticker
	Name         string
	Host         string
	Scheme       string
	Port         int
	sync.Mutex
}

// Refresh is called with the configured TTL time for the given registry
func (r *Registry) Refresh() {

	// Copy the registry information to a new object, and update it
	ur := *r

	// Get the list of repositories
	repoList, err := ur.Registry.Repositories()
	if err != nil {
		utils.Log.WithFields(logrus.Fields{
			"Error": err.Error(),
		}).Error("Failed to retrieve an updated list of repositories for " + ur.URL)
	}
	// Get the repository information
	ur.Repositories = make(map[string]*Repository)
	for _, repoName := range repoList {

		// Build a repository object
		repo := Repository{Name: repoName}

		// Get the list of tags for the repository
		tagList, err := ur.Tags(repoName)
		if err != nil {
			utils.Log.WithFields(logrus.Fields{
				"Error":           err.Error(),
				"Repository Name": repoName,
			}).Error("Failed to retrieve an updated list of tags for " + ur.URL)
		}

		repo.Tags = map[string]*Tag{}
		// Get the manifest for each of the tags
		for _, tagName := range tagList {
			// Use v1 since it has a lot more information
			man, err := ur.Manifest(repoName, tagName)
			if err != nil {
				utils.Log.WithFields(logrus.Fields{
					"Error":           err.Error(),
					"Repository Name": repoName,
					"Tag Name":        tagName,
				}).Error("Failed to retrieve manifest information for " + ur.URL)
			}

			var histories []V1Compatibility
			for _, h := range man.History {
				v1JSON := V1Compatibility{}
				err = json.Unmarshal([]byte(h.V1Compatibility), &v1JSON)
				if err != nil {
					utils.Log.Error(err)
				}
				v1JSON.SizeStr = bytefmt.ByteSize(uint64(v1JSON.Size))

				// Get first 8 characters for the short ID
				v1JSON.IDShort = v1JSON.ID[0:7]

				// Remove shell command
				if len(v1JSON.ContainerConfig.Cmd) > 0 {
					v1JSON.ContainerConfig.CmdClean = strings.Replace(v1JSON.ContainerConfig.Cmd[0], "/bin/sh -c #(nop)", "", -1)
				}
				histories = append(histories, v1JSON)
			}
			repo.Tags[tagName] = &Tag{Name: tagName, V1: man, Histories: histories}
		}
		ur.Repositories[repoName] = &repo
	}
	AllRegistries.Lock()
	AllRegistries.Registries[ur.Name] = &ur
	AllRegistries.Unlock()
}

func (r *Registry) TagCount() int {
	var count int
	for _, repo := range r.Repositories {
		count += len(repo.Tags)
	}
	return count
}

func (r *Registry) Status() string {
	if err := r.Ping(); err != nil {
		return "DOWN"
	}
	return "UP"
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
		Ticker:   time.NewTicker(ttl),
		Host:     host,
		Scheme:   scheme,
		Port:     port,
		Name:     host + ":" + strconv.Itoa(port),
	}
	r.Refresh()

	go func() {
		for range r.Ticker.C {
			r.Refresh()
		}
	}()

	return AllRegistries.Registries[r.Name], nil
}