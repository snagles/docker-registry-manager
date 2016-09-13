package registry

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/Sirupsen/logrus"
	"github.com/stefannaglee/docker-registry-manager/utilities"
)

// Repository contains information on the name and encoded name
type Repository struct {
	Name       string
	Registry   string
	EncodedURI string
	TagCount   int
}

// GetRepositories returns a slice of repositories
func GetRepositories(registryName string) []Repository {
	repos, _ := GetRepositoriesFromRegistry(registryName)
	return repos
}

// GetRepositoriesFromRegistry returns a slice of repositories for this registry name
func GetRepositoriesFromRegistry(registryName string) ([]Repository, error) {

	rs := []Repository{}
	r, err := GetRegistryByName(registryName)
	if err != nil {
		return nil, err
	}

	// Create and execute Get request for the catalog of repositores
	// https://github.com/docker/distribution/blob/master/docs/spec/api.md#catalog
	response, err := http.Get(r.URI() + "/_catalog")
	if err != nil || response.StatusCode != 200 {
		utils.Log.WithFields(logrus.Fields{
			"Registry URL": string(r.URI()),
			"Status Code":  response.StatusCode,
			"Response":     response,
			"Error":        err,
			"Possible Fix": "Check to see if your registry is up, and serving on the correct port with 'docker ps'. ",
		}).Error("Get request to registry failed for the /_catalog endpoint! Is your registry active?")
		return nil, err
	}

	// Close connection
	defer response.Body.Close()

	// Read response into byte body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		utils.Log.WithFields(logrus.Fields{
			"Error": err,
			"Body":  body,
		}).Error("Unable to read response body returned from the registry!")
		return nil, err
	}

	// Unmarshal JSON into the catalog struct containing a slice of repositories
	if unmarshalErr := json.Unmarshal(body, &rs); unmarshalErr != nil {
		utils.Log.WithFields(logrus.Fields{
			"Error":         unmarshalErr,
			"Response Body": string(body),
		}).Error("Unable to unmarshal JSON!")
		return nil, unmarshalErr
	}

	// Escape the the name for the URI
	for _, r := range rs {
		r.EncodedURI = url.QueryEscape(r.Name)
	}

	return rs, err
}
