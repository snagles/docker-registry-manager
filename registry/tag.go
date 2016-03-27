package registry

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

// Tags contains a slice of tags for the given repository
// https://github.com/docker/distribution/blob/master/docs/spec/api.md#listing-image-tags
type Tags struct {
	Name string
	Tags []string
}

// GetTags returns a slice of tags for a given repository and registry
func GetTags(registryName string, repositoryName string) (Tags, error) {

	// Check if the registry is listed as active
	if _, ok := ActiveRegistries[registryName]; !ok {
		return Tags{}, errors.New(registryName + " was not found within the active list of registries.")
	}
	r := ActiveRegistries[registryName]

	// Create and execute Get request
	response, err := http.Get(r.GetURI() + "/" + repositoryName + "/tags/list")
	if err != nil {
		log.WithFields(log.Fields{
			"Registry URL": string(r.GetURI()),
			"Error":        err,
			"Possible Fix": "Check to see if your registry is up, and serving on the correct port with 'docker ps'. ",
		}).Error("Get request to registry failed for the tags endpoint.")
		return Tags{}, err
	}

	// Check Status code
	if response.StatusCode != 200 {
		log.WithFields(log.Fields{
			"Error":       err,
			"Status Code": response.StatusCode,
			"Response":    response,
		}).Error("Did not receive an ok status code!")
		return Tags{}, err
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
		return Tags{}, err
	}
	ts := Tags{}
	// Unmarshal JSON into the tag response struct containing an array of tags
	if err := json.Unmarshal(body, &ts); err != nil {
		log.WithFields(log.Fields{
			"Error":         err,
			"Response Body": string(body),
		}).Error("Unable to unmarshal JSON!")
		return ts, err
	}
	return ts, nil
}
