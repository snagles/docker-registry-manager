package manager

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	client "github.com/snagles/docker-registry-client/registry"
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

	logrus.Info("Refreshing " + r.URL)
	// Get the list of repositories
	repos, err := ur.Registry.Repositories()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err.Error(),
		}).Error("Failed to retrieve an updated list of repositories for " + ur.URL)
	}
	// Get the repository information
	ur.Repositories = make(map[string]*Repository)
	for _, repoName := range repos {

		// Get the list of tags for the repository
		tags, err := ur.Tags(repoName)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"Error":           err.Error(),
				"Repository Name": repoName,
			}).Error("Failed to retrieve an updated list of tags for " + ur.URL)
			continue
		}

		repo := Repository{Name: repoName, Tags: make(map[string]*Tag)}
		// Get the manifest for each of the tags
		for _, tagName := range tags {

			// Using v2 required getting the manifest then retrieving the blob
			// for the config digest
			man, err := ur.Manifest(repoName, tagName)
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
				continue
			}
			var v1 V1Compatibility
			err = json.Unmarshal(v1Bytes, &v1)
			if err != nil {
				logrus.Error(err)
				continue
			}

			// add the pointer for the history to its layer using its index
			layerIndex := 0
			for i, history := range v1.History {
				if !history.EmptyLayer {
					v1.History[i].ManifestLayer = &man.Layers[layerIndex]
					layerIndex++
				}
				sh := strings.Split(history.CreatedBy, "/bin/sh -c")
				if len(sh) > 1 {
					v1.History[i].ShellType = "/bin/sh -c"
					commands := strings.SplitAfter(sh[1], "&&")
					for _, cmd := range commands {
						v1.History[i].Commands = append(v1.History[i].Commands, Command{Cmd: cmd, Keywords: Keywords(cmd)})
					}
				}
			}

			// Get the tag size information
			size, err := ur.TagSize(repoName, tagName)
			if err != nil {
				logrus.Error(err)
			}

			repo.Tags[tagName] = &Tag{Name: tagName, V1Compatibility: &v1, Size: int64(size), DeserializedManifest: man}
		}
		ur.Repositories[repoName] = &repo
	}
	AllRegistries.Lock()
	AllRegistries.Registries[ur.Name] = &ur
	AllRegistries.Unlock()
}

func (r *Registry) TagCount() (count int) {
	for _, repo := range r.Repositories {
		count += len(repo.Tags)
	}
	return count
}

func (r *Registry) LayerCount() int {
	layerDigests := make(map[string]struct{})
	for _, repo := range r.Repositories {
		for _, tag := range repo.Tags {
			for _, layer := range tag.Layers {
				layerDigests[layer.Digest.String()] = struct{}{}
			}
		}
	}
	return len(layerDigests)
}

func (r *Registry) Pushes() (pushes int) {
	AllEvents.Lock()
	defer AllEvents.Unlock()
	if _, ok := AllEvents.Events[r.Name]; !ok {
		return 0
	}

	for _, e := range AllEvents.Events[r.Name] {
		// TODO: really need to find a better way to exclude the managers queries
		if e.Action == "push" && e.Request.Useragent != "Go-http-client/1.1" && e.Request.Method != "HEAD" {
			pushes++
		}
	}
	return pushes
}

func (r *Registry) Pulls() (pulls int) {
	AllEvents.Lock()
	defer AllEvents.Unlock()
	if _, ok := AllEvents.Events[r.Name]; !ok {
		return 0
	}

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

// AddRegistry adds the new registry for viewing in the interface and sets up
// the go routine for automatic refreshes
func AddRegistry(scheme, host, user, password string, port int, ttl time.Duration, skipTLS bool) (*Registry, error) {
	switch {
	case scheme == "":
		return nil, errors.New("Invalid scheme: " + scheme)
	case host == "":
		return nil, errors.New("Invalid host: " + host)
	case port == 0:
		return nil, errors.New("Invalid port: " + strconv.Itoa(port))
	}

	url := fmt.Sprintf("%s://%s:%v", scheme, host, port)
	var hub *client.Registry
	var err error
	if skipTLS {
		hub, err = client.NewInsecure(url, user, password)
		if err != nil {
			logrus.Error("Failed to connect to unvalidated TLS registry: " + err.Error())
			return nil, err
		}
	} else {
		hub, err = client.New(url, user, password)
		if err != nil {
			logrus.Error("Failed to connect to validated registry: " + err.Error())
			return nil, err
		}
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
