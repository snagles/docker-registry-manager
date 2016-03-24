package registry

import (
	"net"
	"net/http"
	"net/url"

	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql" // need to initialize mysql before making a connection
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
}

// GetURI returns the full url path for communicating with this registry
func (r *Registry) GetURI() string {
	return r.Scheme + "://" + r.Name + ":" + r.Port + "/v2"
}

// TestRegistryStatus takes in a registry URL and checks for communication errors
//
// Create and execute basic GET request to test if each registry can be reached
// To determine registry status we test the base registry route of /v2/ and check
// the HTTP response code for a 200 response (200 is a successful request)
func TestRegistryStatus(registryURI string) error {

	// Parse the registry string into our Registry type
	registry := ParseRegistry(registryURI)
	// Notify of initial attempt
	log.WithFields(log.Fields{
		"Registry URI": registryURI,
	}).Info("Connecting to registry...")

	// Create and execute a plain get request and check the http status code
	response, err := http.Get(registry.GetURI())
	if err != nil {
		// Notify of error
		log.WithFields(log.Fields{
			"Registry URLs": registry,
			"Error":         err,
			"HTTP Response": response,
			"Possible Fix":  "Check to see if your registry is up, and serving on the correct port with 'docker ps'.",
		}).Fatal("Get request to registry timed out/failed! Is the URL correct, and is the registry active?")

		return err
	} else if response.StatusCode != 200 {
		// Notify of error
		log.WithFields(log.Fields{
			"Registry URLs": registry,
			"HTTP Response": response.StatusCode,
			"Possible Fix":  "Check to see if your registry is up, and serving on the correct port with 'docker ps'.",
		}).Fatal("Get request to registry failed! Is the URL correct, and is the registry active?")
	}

	// Notify of success
	log.WithFields(log.Fields{
		"Registry Information": registry,
		"Registry URI":         registry.GetURI(),
	}).Info("Successfully connected to registry and added to list of active registries!")

	// Add the registry to the map of active registries
	ActiveRegistries[registry.Name] = registry

	return err
}

// ParseRegistry takes in a registry URI string and converts it into a registry object
func ParseRegistry(registryURI string) Registry {

	// Parse the URL and get the scheme
	// e.g https, http, etc.
	u, err := url.Parse(registryURI)
	if err != nil {
		log.Error(err)
	}

	// Get the host and port
	// e.g test.domain.com and 5000, etc.
	host, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		log.Error(err)
	}

	// Lookup the ip for the passed host
	// Using the host name try looking up the IP for informational purposes
	ip, err := net.LookupHost(host)
	if err != nil {
		log.Error(err)
	}

	// Create a new registry type using all of the information queried above
	r := Registry{
		Name:    host,
		Port:    port,
		IP:      ip[0],
		Scheme:  u.Scheme,
		Version: "v2",
	}

	// Return the newly created registry type
	return r
}
