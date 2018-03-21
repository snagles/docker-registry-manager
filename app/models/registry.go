package manager

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	manifestV2 "github.com/docker/distribution/manifest/schema2"
	"github.com/sirupsen/logrus"
	client "github.com/snagles/docker-registry-client/registry"
	"github.com/spf13/viper"
)

const (
	StatusUp   = "UP"
	StatusDown = "DOWN"
)

// AllRegistries contains a list of added registries using their hostnames
// access granted via mutex locks/unlocks
var AllRegistries Registries

func init() {
	AllRegistries.Registries = map[string]*Registry{}
}

// Registries contains a map of all active registries identified by their name, locked when necessary
type Registries struct {
	Registries map[string]*Registry
	*viper.Viper
	sync.Mutex
}

type registriesConfig struct {
	URL                  string
	Port                 int
	Username             string
	Password             string
	SkipTLS              bool   `mapstructure:"skip-tls-validation" yaml:"skip-tls-validation"`
	RefreshRate          string `mapstructure:"refresh-rate" yaml:"refresh-rate"`
	DockerhubIntegration bool   `mapstructure:"dockerhub-integration" yaml:"dockerhub-integration"`
}

// AddRegistry adds the registry to the all registries map its details
func (rs *Registries) AddRegistry(scheme, host, name, user, password string, port int, ttl time.Duration, skipTLS, dockerhubIntegration bool) error {
	r, err := NewRegistry(scheme, host, name, user, password, port, ttl, skipTLS, dockerhubIntegration)
	if err != nil {
		return err
	}
	AllRegistries.Lock()
	AllRegistries.Registries[r.Name] = r
	AllRegistries.Unlock()
	logrus.Infof("Added new registry: %s", name)
	return nil
}

// RemoveRegistry removes the registry from the all registries map using its name
func (rs *Registries) RemoveRegistry(name string) {
	AllRegistries.Lock()
	AllRegistries.Registries[name].Ticker.Stop()
	AllRegistries.Registries[name].StopRefresh <- false
	delete(AllRegistries.Registries, name)
	AllRegistries.Unlock()
	logrus.Infof("Removed registry: %s", name)
}

// LoadConfig adds the registries parsed from the passed yaml file
func (rs *Registries) LoadConfig(registriesFile string) {
	if rs.Viper == nil {
		rs.Viper = viper.New()
	}

	// If the registries path is not passed use the default project dir
	if registriesFile != "" {
		rs.AddConfigPath(path.Dir(registriesFile))
		base := path.Base(registriesFile)
		ext := path.Ext(registriesFile)
		rs.SetConfigName(base[0 : len(base)-len(ext)])
		logrus.Infof("Using registries located in %s with file name %s", path.Dir(registriesFile), base[0:len(base)-len(ext)])
	} else {
		rs.SetConfigName("registries")
		var root string
		_, run, _, ok := runtime.Caller(0)
		if ok {
			root = filepath.Dir(run)
			rs.AddConfigPath(root)
		} else {
			logrus.Fatalf("Failed to get runtime caller for parser")
		}
		logrus.Infof("Using registries located in %s with file name %s", root, "registries.yml")
	}

	config := make(map[string]map[string]registriesConfig)

	if err := rs.ReadInConfig(); err != nil {
		logrus.Fatalf("Failed to read in registries file: %s", err)
	}

	if err := rs.Unmarshal(&config); err != nil {
		logrus.Fatalf("Unable to unmarshal registries file: %s", err)
	}

	// overwrite the entries with the updated information
	for name, r := range config["registries"] {
		if r.URL != "" {
			url, err := url.Parse(r.URL)
			if err != nil {
				logrus.Fatalf("Failed to parse registry from the passed url (%s): %s", r.URL, err)
			}
			duration, err := time.ParseDuration(r.RefreshRate)
			if err != nil {
				logrus.Fatalf("Failed to add registry (%s), invalid duration: %s", r.URL, err)
			}
			if err := AllRegistries.AddRegistry(url.Scheme, url.Hostname(), name, r.Username, r.Password, r.Port, duration, r.SkipTLS, r.DockerhubIntegration); err != nil {
				logrus.Fatalf("Failed to add registry (%s): %s", r.URL, err)
			}
		}
	}
}

// WriteConfig builds the config and writes from the map of registries
func (rs *Registries) WriteConfig() error {
	config := make(map[string]registriesConfig)
	for name, r := range AllRegistries.Registries {
		config[name] = registriesConfig{
			URL:                  r.URL,
			Port:                 r.Port,
			Username:             r.Username,
			Password:             r.Password,
			SkipTLS:              r.SkipTLS,
			RefreshRate:          r.TTL.String(),
			DockerhubIntegration: r.DockerhubIntegration,
		}
	}

	rs.Set("registries", config)
	logrus.Info("Writing config with new/removed registries")
	return rs.Viper.WriteConfig()
}

// Registry contains all information about the registry and its metadata
type Registry struct {
	*client.Registry
	Repositories         map[string]*Repository
	TTL                  time.Duration
	Ticker               *time.Ticker
	StopRefresh          chan (bool)
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

// Refresh is called with the configured TTL time for the given registry
func (r Registry) Refresh() {
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
	// Add history at most every 15 minutes
	if len(r.History) == 0 || n.Sub(r.History[len(r.History)-1].Time) >= (15*time.Minute) {
		r.History = append(r.History, RegistryHistory{
			Repositories: len(r.Repositories),
			Tags:         r.TagCount(),
			Layers:       r.LayerCount(),
			Time:         n,
		})
	}

	// purge everything older than 3 days
	for i, h := range r.History {
		if h.Time.Before(n.AddDate(0, -3, 0)) {
			r.History = append(r.History[:i], r.History[i+1:]...)
		}
	}
	r.LastRefresh = time.Now().UTC()
	AllRegistries.Registries[r.Name] = &r
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
		Ticker:               time.NewTicker(ttl),
		StopRefresh:          make(chan bool, 1),
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

	r.Refresh()
	go func() {
		for {
			select {
			case <-r.Ticker.C:
				r.Refresh()
			case <-r.StopRefresh:
				return
			}
		}
	}()
	return &r, nil
}
