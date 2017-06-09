package manager

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	client "github.com/heroku/docker-registry-client/registry"
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
	Version      string
	Port         int
	sync.Mutex
}

func (r *Registry) IP() string {
	ip, _ := net.LookupHost(r.Host)
	if len(ip) > 0 {
		return ip[0]
	}
	return ""
}

// Refresh is called with the configured TTL time for the given registry
func (r *Registry) Refresh() {

	// Copy the registry information to a new object, and update it
	ur := *r

	// Get the list of repositories
	repoList, err := ur.Registry.Repositories()
	if err != nil {
		logrus.WithFields(logrus.Fields{
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
			logrus.WithFields(logrus.Fields{
				"Error":           err.Error(),
				"Repository Name": repoName,
			}).Error("Failed to retrieve an updated list of tags for " + ur.URL)
			return
		}

		repo.Tags = map[string]*Tag{}
		// Get the manifest for each of the tags
		for _, tagName := range tagList {

			// Using v2 required getting the manifest then retrieving the blob
			// for the config digest
			man, err := ur.ManifestV2(repoName, tagName)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"Error":           err.Error(),
					"Repository Name": repoName,
					"Tag Name":        tagName,
				}).Error("Failed to retrieve manifest information for " + ur.URL)
				continue
			}

			// Get the v1 config information
			v1Bytes, err := ur.ManifestMetadata(repoName, man.Config.Digest)
			if err != nil {
				logrus.Error(err)
				return
			}
			var v1 V1Compatibility
			err = json.Unmarshal(v1Bytes, &v1)
			if err != nil {
				logrus.Error(err)
				return
			}

			// add the pointer for the history to its layer
			layerIndex := 0
			for i, history := range v1.History {
				if !history.EmptyLayer {
					v1.History[i].ManifestLayer = &man.Layers[layerIndex]
				} else {
					layerIndex++
				}
			}

			// Get the tag size information
			size, err := ur.TagSize(repoName, tagName)
			if err != nil {
				logrus.Error(err)
				return
			}

			repo.Tags[tagName] = &Tag{Name: tagName, V1Compatibility: &v1, Size: int64(size), DeserializedManifest: man}
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

func (r *Registry) LayerCount() int {
	var count int
	for _, repo := range r.Repositories {
		for _, tag := range repo.Tags {
			count += tag.LayerCount()
		}
	}
	return count
}

func (r *Registry) Pushes() int {
	AllEvents.Lock()
	defer AllEvents.Unlock()
	if _, ok := AllEvents.Events[r.Name]; !ok {
		return 0
	}

	var pushes int
	for _, e := range AllEvents.Events[r.Name] {
		// TODO: really need to find a better way to exclude the managers queries
		if e.Action == "push" && e.Request.Useragent != "Go-http-client/1.1" && e.Request.Method != "HEAD" {
			pushes++
		}
	}
	return pushes
}

func (r *Registry) Pulls() int {
	AllEvents.Lock()
	defer AllEvents.Unlock()
	if _, ok := AllEvents.Events[r.Name]; !ok {
		return 0
	}

	var pulls int
	for _, e := range AllEvents.Events[r.Name] {
		// exclude heads since thats the method the manager uses for getting meta info
		// TODO: really need to find a better way to exclude the managers queries
		if e.Action == "pull" && e.Request.Useragent != "Go-http-client/1.1" && e.Request.Method != "HEAD" {
			pulls++
		}
	}
	return pulls
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
		Version:  "v2",
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
