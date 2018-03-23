package manager

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	manifestV2 "github.com/docker/distribution/manifest/schema2"
	"github.com/sirupsen/logrus"
	client "github.com/snagles/docker-registry-client/registry"
)

const (
	StatusUp   = "UP"
	StatusDown = "DOWN"
)

// Registry contains all information about the registry and its metadata
type Registry struct {
	*client.Registry
	Repositories         map[string]*Repository
	TTL                  time.Duration
	Name                 string
	Username             string
	Password             string
	Host                 string
	Scheme               string
	Version              string
	Port                 int
	DockerhubIntegration bool
	SkipTLS              bool
	LastRefresh          time.Time
	status               string
	ip                   string
	History              []RegistryHistory
}

func (r *Registry) HistoryTimes() []time.Time {
	var hs []time.Time
	for i := range r.History {
		hs = append(hs, r.History[i].Time)
	}
	return hs
}

func (r *Registry) HistoryRepos() []int {
	var hs []int
	for i := range r.History {
		hs = append(hs, r.History[i].Repositories)
	}
	return hs
}

func (r *Registry) HistoryLayers() []int {
	var hs []int
	for i := range r.History {
		hs = append(hs, r.History[i].Layers)
	}
	return hs
}

func (r *Registry) HistoryTags() []int {
	var hs []int
	for i := range r.History {
		hs = append(hs, r.History[i].Tags)
	}
	return hs
}

// RegistryHistory maintains a list of data points at a regular interval for plotting in the UI
type RegistryHistory struct {
	Repositories int
	Layers       int
	Tags         int
	Time         time.Time
}

// IP returns the ip as a string
func (r *Registry) IP() string {
	return r.ip
}

// Update is called with the configured TTL time for the given registry
func (r Registry) Update() Registry {
	old := r
	err := r.Ping()
	if err != nil {
		r.status = StatusDown
	} else {
		r.status = StatusUp
	}

	ip, _ := net.LookupHost(r.Host)
	if len(ip) > 0 {
		r.ip = ip[0]
	}

	logrus.Info("Refreshing " + r.Name)
	// Get the list of repositories
	repos, err := r.Registry.Repositories()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err.Error(),
		}).Error("Failed to retrieve an updated list of repositories for " + r.URL)
	}
	// Get the repository information
	r.Repositories = make(map[string]*Repository)
	for _, repoName := range repos {

		// Get the list of tags for the repository
		tags, err := r.Tags(repoName)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"Error":           err.Error(),
				"Repository Name": repoName,
			}).Error("Failed to retrieve an updated list of tags for " + r.URL)
			continue
		}

		repo := Repository{Name: repoName, Tags: make(map[string]*Tag)}
		// Get the manifest for each of the tags
		for _, tagName := range tags {

			// Using v2 required getting the manifest then retrieving the blob
			// for the config digest
			man, err := r.Manifest(repoName, tagName)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"Error":           err.Error(),
					"Repository Name": repoName,
					"Tag Name":        tagName,
				}).Error("Failed to retrieve manifest information for " + r.URL)
				continue
			}

			// Get the v1 config information
			v1Bytes, err := r.ManifestMetadata(repoName, man.Config.Digest)
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

			// Get the tag size information from the manifest layers
			size, err := r.CalculateTagSize(man)
			if err != nil {
				logrus.Error(err)
			}

			repo.Tags[tagName] = &Tag{Name: tagName, V1Compatibility: &v1, Size: int64(size), DeserializedManifest: man}
		}
		r.Repositories[repoName] = &repo
	}

	n := time.Now().UTC()
	if len(r.History) == 0 || old.LayerCount() != r.LayerCount() || old.TagCount() != r.TagCount() || len(old.Repositories) != len(r.Repositories) {
		r.History = append(r.History, RegistryHistory{
			Repositories: len(r.Repositories),
			Tags:         r.TagCount(),
			Layers:       r.LayerCount(),
			Time:         n,
		})
	}

	// purge everything older than 5 days
	for i, h := range r.History {
		if h.Time.Before(n.AddDate(0, -5, 0)) {
			r.History = append(r.History[:i], r.History[i+1:]...)
		}
	}
	r.LastRefresh = time.Now().UTC()
	return r
}

// CalculateTagSize returns the total number of tags across all repositories
func (r *Registry) CalculateTagSize(deserialized *manifestV2.DeserializedManifest) (size int64, err error) {
	size = int64(0)
	for _, layer := range deserialized.Layers {
		size += layer.Size
	}
	return size, nil
}

// TagCount returns the total number of tags across all repositories
func (r *Registry) TagCount() (count int) {
	for _, repo := range r.Repositories {
		count += len(repo.Tags)
	}
	return count
}

// LayerCount returns the total number of layers across all repositories
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

// Pushes returns the number of pushes recorded by passing the forwarded registry events
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

// Pulls returns the number of pulls recorded by passing the forwarded registry events
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

// Status returns the text representation of whether the registry is reachable
func (r *Registry) Status() string {
	return r.status
}

// NewRegistry adds the new registry for viewing in the interface and sets up
// the go routine for automatic refreshes
func NewRegistry(scheme, host, name, user, password string, port int, ttl time.Duration, skipTLS, dockerhubIntegration bool) (*Registry, error) {
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
		Registry:             hub,
		TTL:                  ttl,
		Host:                 host,
		Scheme:               scheme,
		Username:             user,
		Password:             password,
		Port:                 port,
		Version:              "v2",
		Name:                 name,
		SkipTLS:              skipTLS,
		DockerhubIntegration: dockerhubIntegration,
	}
	return &r, nil
}
