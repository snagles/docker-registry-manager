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

	repos, err := client.GetRepositories(uri)
	if err != nil {
		return err
	}

	// Use a layer map to de-duplicate shared layers
	layerMap := make(map[string]int64, 0)

	// Get the tags and image information
	for _, repoName := range repos {
		r.Repositories = append(r.Repositories, Repository{Name: repoName})

		tagList, err := client.GetTags(uri, repoName)
		if err == nil {
			r.TagCount = len(tagList)
			for _, tagName := range tagList {
				img, _ := client.GetImage(uri, repoName, tagName)
				for _, layer := range img.FsLayers {
					layerMap[layer.BlobSum] = layer.Size
				}
			}
		}
	}

	// Total the size of the layers
	for _, size := range layerMap {
		r.RepoTotalSize += size
	}
	r.RepoTotalSizeStr = bytefmt.ByteSize(uint64(r.RepoTotalSize))

	Registries[r.Name] = &r

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

type Repository struct {
	Name string
	Tags []Tag
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
