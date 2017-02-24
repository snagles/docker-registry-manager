package registry

import (
	"net"
	"time"

	"fmt"

	"github.com/snagles/docker-registry-manager/models/client"
	"github.com/snagles/docker-registry-manager/utilities"

	"code.cloudfoundry.org/bytefmt"
)

// Registries contains a map of all active registries identified by their name
var Registries map[string]*registry

func init() {
	Registries = map[string]*registry{}
}

func refresh() {
	for {
		utils.Log.Debug("Start refreshing")
		for name, registry := range Registries {
			utils.Log.Debug("Refreshing " + name)
			registry.Refresh()
			utils.Log.Debug("Refreshed " + name)
		}
		utils.Log.Debug("Refreshing is done")
		time.Sleep(45 * time.Second)
	}
}

func AddRegistry(scheme, host string, port int) error {
	r, err := NewRegistry(scheme, host, port)
	if err != nil {
		return err
	}

	err = client.HealthCheck(r.URI())
	if err != nil {
		return err
	}
	r.Refresh()
	go refresh()

	return nil
}

// Create a new registry from URI parts
func NewRegistry(scheme, host string, port int) (registry, error) {
	var err error
	r := registry{
		Scheme:  scheme,
		Host:    host,
		Port:    port,
		Version: 2,
	}

	// Lookup the ip for the passed host
	ip, err := net.LookupHost(host)
	if err != nil {
		utils.Log.Error(err)
		return registry{}, err
	}
	if len(ip) > 0 {
		r.IP = ip[0]
	}

	return r, nil
}

// Registry contains all identifying information for communicating with a registry
type registry struct {
	Host         string
	IP           string
	Scheme       string
	Port         int
	Version      int
	Repositories map[string]*Repository
	Status       int
}

// URI returns the full url path for communicating with this registry
func (r *registry) URI() string {
	return fmt.Sprintf("%s://%s:%d/v%d", r.Scheme, r.Host, r.Port, r.Version)
}

func (r *registry) TagCount() int {
	var count int
	for _, repo := range r.Repositories {
		count += len(repo.Tags)
	}
	return count
}

func (r *registry) DiskSize() string {
	registryLayerMap := map[string]bool{}
	var size int64

	for _, repo := range r.Repositories {
		for _, tag := range repo.Tags {
			for _, layer := range tag.Image.FsLayers {
				if _, ok := registryLayerMap[layer.BlobSum]; !ok {
					registryLayerMap[layer.BlobSum] = true
					size += layer.Size
				}
			}
		}
	}
	return bytefmt.ByteSize(uint64(size))
}

func (r registry) Refresh() {

	// Get the list of repositories on this registry
	repoList, _ := client.GetRepositories(r.URI())

	// Get the repository information
	r.Repositories = map[string]*Repository{}
	for _, repoName := range repoList {

		// Build a repository object
		repo := Repository{Name: repoName}

		tagList, _ := client.GetTags(r.URI(), repoName)
		repo.Tags = map[string]*Tag{}
		for _, tagName := range tagList {
			tag := Tag{Name: tagName}
			img, _ := client.GetImage(r.URI(), repoName, tagName)
			tag.Image = Image{img}

			// Add the tag to the repository
			repo.Tags[tagName] = &tag
		}
		r.Repositories[repoName] = &repo
	}
	Registries[r.Host] = &r
}

type Repository struct {
	Name string
	Tags map[string]*Tag
}

func (r *Repository) LastModified() time.Time {
	var lastModified time.Time
	for _, tag := range r.Tags {
		for _, history := range tag.Image.History {
			if history.V1Compatibility.Created.After(lastModified) {
				lastModified = history.V1Compatibility.Created
			}
		}
	}
	return lastModified
}

func (r *Repository) LastModifiedTimeAgo() string {
	lastModified := r.LastModified()
	return utils.TimeAgo(lastModified)
}

type Tag struct {
	Image

	ID   string
	Name string
}

func (t *Tag) LastModified() time.Time {
	var lastModified time.Time
	for _, history := range t.Image.History {
		if history.V1Compatibility.Created.After(lastModified) {
			lastModified = history.V1Compatibility.Created
		}
	}
	return lastModified
}

func (t *Tag) LastModifiedTimeAgo() string {
	lastModified := t.LastModified()
	return utils.TimeAgo(lastModified)
}

func (t *Tag) Size() string {
	var size int64
	for _, layer := range t.Image.FsLayers {
		size += layer.Size
	}
	return bytefmt.ByteSize(uint64(size))
}

func (t *Tag) LayerCount() int {
	return len(t.Image.FsLayers)
}

type Image struct {
	client.Image
}

func (i *Image) LastModified() time.Time {
	var lastModified time.Time
	for _, history := range i.History {
		if history.V1Compatibility.Created.After(lastModified) {
			lastModified = history.V1Compatibility.Created
		}
	}
	return lastModified
}

func (i *Image) LastModifiedTimeAgo() string {
	lastModified := i.LastModified()
	return utils.TimeAgo(lastModified)
}

func (i *Image) Size() string {
	var size int64
	for _, layer := range i.FsLayers {
		size += layer.Size
	}
	return bytefmt.ByteSize(uint64(size))
}

func (i *Image) LayerCount() int {
	return len(i.FsLayers)
}
