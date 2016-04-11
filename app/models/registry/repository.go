package registry

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/Sirupsen/logrus"
	"github.com/stefannaglee/docker-registry-manager/app/utilities"
)

// RepositoriesList contains a slice of all repositories
type RepositoriesList struct {
	Repositories []string
}

// Repository contains information on the name and encoded name
type Repository struct {
	Name       string
	EncodedURI string
}

// GetRepositories returns a slice of repositories with their names and encoded names
func GetRepositories(registryName string) []Repository {
	cleanedRepos := []Repository{}
	repos, _ := GetRepositoriesFromRegistry(registryName)
	for _, value := range repos.Repositories {
		r := Repository{}
		r.EncodedURI = url.QueryEscape(value)
		r.Name = value
		cleanedRepos = append(cleanedRepos, r)
	}
	return cleanedRepos
}

// GetRepositoriesFromRegistry returns a slice of repositories for this registry name
func GetRepositoriesFromRegistry(registryName string) (RepositoriesList, error) {

	// Check if the registry is listed as active
	if _, ok := ActiveRegistries[registryName]; !ok {
		return RepositoriesList{}, errors.New(registryName + " was not found within the active list of registries.")
	}
	r := ActiveRegistries[registryName]

	// Create and execute Get request for the catalog of repositores
	// https://github.com/docker/distribution/blob/master/docs/spec/api.md#catalog
	response, err := http.Get(r.GetURI() + "/_catalog")
	if err != nil {
		utils.Log.WithFields(logrus.Fields{
			"Registry URL": string(r.GetURI()),
			"Error":        err,
			"Possible Fix": "Check to see if your registry is up, and serving on the correct port with 'docker ps'. ",
		}).Error("Get request to registry failed for the /_catalog endpoint! Is your registry active?")
	}

	// Check Status code
	if response.StatusCode != 200 {
		utils.Log.WithFields(logrus.Fields{
			"Error":       err,
			"Status Code": response.StatusCode,
			"Response":    response,
		}).Error("Did not receive an ok status code!")
		return RepositoriesList{}, err
	}

	// Close connection
	defer response.Body.Close()

	// Read response into byte body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		utils.Log.WithFields(logrus.Fields{
			"Error": err,
			"Body":  body,
		}).Error("Unable to read response into body!")
	}

	rs := RepositoriesList{}
	// Unmarshal JSON into the catalog struct containing a slice of repositories
	if err := json.Unmarshal(body, &rs); err != nil {
		utils.Log.WithFields(logrus.Fields{
			"Error":         err,
			"Response Body": string(body),
		}).Error("Unable to unmarshal JSON!")
	}

	return rs, err
}
