package registry

import (
	"net"
	"net/url"
	"time"

	"github.com/pivotal-golang/bytefmt"
	"github.com/snagles/docker-registry-manager/models/client"
	"github.com/snagles/docker-registry-manager/utilities"
)

// Registries contains a map of all active registries identified by their name
var Registries map[string]*Registry

func init() {
	Registries = make(map[string]*Registry, 0)
}

func AddRegistry(uri string) error {
	r, err := ParseRegistry(uri)
	if err != nil {
		return err
	}

	r.Refresh()

	return nil
}

// Registry contains all identifying information for communicating with a registry
type Registry struct {
	Name             string
	IP               string
	Scheme           string
	Port             string
	Version          string
	Repositories     []Repository
	TagCount         int
	Status           string
	RepoTotalSize    int64
	RepoTotalSizeStr string
}

// URI returns the full url path for communicating with this registry
func (r *Registry) URI() string {
	return r.Scheme + "://" + r.Name + ":" + r.Port + "/v2"
}

func (r *Registry) Refresh() {

	// Get the lsit of repositories on this registry
	repoList, _ := client.GetRepositories(r.URI())

	// Use a layer map to de-duplicate shared layers across the registry
	registryLayerMap := make(map[string]int64, 0)

	// Get the repository information
	for _, repoName := range repoList {

		// Build a repository object
		repo := Repository{Name: repoName}

		tagList, _ := client.GetTags(r.URI(), repoName)
		for _, tagName := range tagList {
			tag := Tag{Name: tagName}
			tag.Image, _ = client.GetImage(r.URI(), repoName, tagName)
			for _, layer := range tag.Image.FsLayers {
				registryLayerMap[layer.BlobSum] = layer.Size
			}

			// Add the tag to the repository
			repo.Tags = append(repo.Tags, tag)
		}
		r.Repositories = append(r.Repositories, repo)
		r.TagCount += len(tagList)
	}

	// Total the size of the layers
	for _, size := range registryLayerMap {
		r.RepoTotalSize += size
	}
	r.RepoTotalSizeStr = bytefmt.ByteSize(uint64(r.RepoTotalSize))

	Registries[r.Name] = r
}

type Repository struct {
	Name string
	Tags []Tag
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
	Image client.Image

	ID              string
	Name            string
	UpdatedTime     time.Time
	UpdatedTimeUnix int64
	TimeAgo         string
	Layers          int
	Size            string
	SizeInt         int64
}

// ParseRegistry takes in a registry URI string and converts it into a registry object
func ParseRegistry(registry string) (Registry, error) {

	// only supports v2 currently
	r := Registry{
		Version: "v2",
	}

	// Parse the URL and get the scheme
	// e.g https, http, etc.
	u, err := url.Parse(registry)
	if err != nil {
		utils.Log.Error(err)
		return r, err
	}

	// Set scheme
	r.Scheme = u.Scheme

	// Get the host and port
	// e.g test.domain.com and 5000, etc.
	host, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		utils.Log.Error(err)
		return r, err
	}

	// Set name and port
	r.Name = host
	r.Port = port

	// Lookup the ip for the passed host
	ip, err := net.LookupHost(host)
	if err != nil {
		utils.Log.Error(err)
		return r, err
	}
	if len(ip) > 0 {
		r.IP = ip[0]
	}

	return r, err
}
