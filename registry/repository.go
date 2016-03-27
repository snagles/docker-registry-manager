package registry

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

// Repositories contains a slice of all repositories
type Repositories struct {
	Repositories []string
}

// GetRepositories returns a slice of repositories for this registry name
func GetRepositories(registryName string) (Repositories, error) {

	// Check if the registry is listed as active
	if _, ok := ActiveRegistries[registryName]; !ok {
		return Repositories{}, errors.New(registryName + " was not found within the active list of registries.")
	}
	r := ActiveRegistries[registryName]

	// Create and execute Get request for the catalog of repositores
	// https://github.com/docker/distribution/blob/master/docs/spec/api.md#catalog
	response, err := http.Get(r.GetURI() + "/_catalog")
	if err != nil {
		log.WithFields(log.Fields{
			"Registry URL": string(r.GetURI()),
			"Error":        err,
			"Possible Fix": "Check to see if your registry is up, and serving on the correct port with 'docker ps'. ",
		}).Error("Get request to registry failed for the /_catalog endpoint! Is your registry active?")
	}

	// Check Status code
	if response.StatusCode != 200 {
		log.WithFields(log.Fields{
			"Error":       err,
			"Status Code": response.StatusCode,
			"Response":    response,
		}).Error("Did not receive an ok status code!")
		return Repositories{}, err
	}

	// Close connection
	defer response.Body.Close()

	// Read response into byte body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
			"Body":  body,
		}).Error("Unable to read response into body!")
	}

	rs := Repositories{}
	// Unmarshal JSON into the catalog struct containing a slice of repositories
	if err := json.Unmarshal(body, &rs); err != nil {
		log.WithFields(log.Fields{
			"Error":         err,
			"Response Body": string(body),
		}).Error("Unable to unmarshal JSON!")
	}

	return rs, err
}
