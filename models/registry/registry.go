package registry

import (
	"net"
	"net/http"
	"net/url"

	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql" // need to initialize mysql before making a connection
	"github.com/stefannaglee/docker-registry-manager/utilities"
)

// ActiveRegistries contains a map of all active registries identified by their name
var ActiveRegistries map[string]Registry

func init() {
	// Create the active registries map
	ActiveRegistries = make(map[string]Registry, 0)
}

// Registry contains all identifying information for communicating with a registry
type Registry struct {
	Name    string
	IP      string
	Scheme  string
	Port    string
	Version string

	Status           string
	RepoCount        int
	TagCount         int
	RepoTotalSize    int64
	RepoTotalSizeStr string
}

// GetURI returns the full url path for communicating with this registry
func (r *Registry) GetURI() string {
	return r.Scheme + "://" + r.Name + ":" + r.Port + "/v2"
}

// UpdateRegistryStatus takes in a registry URL and checks for communication errors
//
// Create and execute basic GET request to test if each registry can be reached
// To determine registry status we test the base registry route of /v2/ and check
// the HTTP response code for a 200 response (200 is a successful request)
func (r *Registry) UpdateRegistryStatus() error {

	// Parse the registry string into our Registry type
	utils.Log.WithFields(logrus.Fields{
		"Registry URI": r.GetURI(),
	}).Info("Connecting to registry...")

	// Create and execute a plain get request and check the http status code
	response, err := http.Get(r.GetURI())
	if err != nil {
		// Notify of error
		utils.Log.WithFields(logrus.Fields{
			"Registry URLs": r,
			"Error":         err,
			"HTTP Response": response,
			"Possible Fix":  "Check to see if your registry is up, and serving on the correct port with 'docker ps'.",
		}).Error("Get request to registry timed out/failed! Is the URL correct, and is the registry active?")
		r.Status = "unavailable"

		return err
	} else if response.StatusCode != 200 {
		// Notify of error
		utils.Log.WithFields(logrus.Fields{
			"Registry URLs": r,
			"HTTP Response": response.StatusCode,
			"Possible Fix":  "Check to see if your registry is up, and serving on the correct port with 'docker ps'.",
		}).Error("Get request to registry failed! Is the URL correct, and is the registry active?")
		r.Status = "unavailable"
	}

	// Notify of success
	utils.Log.WithFields(logrus.Fields{
		"Registry Information": r,
		"Registry URI":         r.GetURI(),
	}).Info("Successfully connected to registry and added to list of active registries!")

	r.Status = "available"

	return err
}

// AddRegistry adds the registry to the map of active registries
func (r *Registry) AddRegistry() {
	ActiveRegistries[r.Name] = *r
}

// ParseRegistry takes in a registry URI string and converts it into a registry object
func ParseRegistry(registryURI string) (Registry, error) {

	// Create an empty Registry
	r := Registry{
		Version: "v2",
	}

	// Parse the URL and get the scheme
	// e.g https, http, etc.
	u, err := url.Parse(registryURI)
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
	// Using the host name try looking up the IP for informational purposes
	ip, err := net.LookupHost(host)
	if err != nil {
		utils.Log.Error(err)
		// We do not need to return an error since we don't "need" the IP of the host
	}
	// Set IP if we have it
	if ip != nil {
		r.IP = ip[0]
	}

	// Return the newly created registry type
	return r, err
}
