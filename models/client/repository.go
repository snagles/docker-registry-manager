package client

import (
	"encoding/json"

	"github.com/DemonVex/docker-registry-manager/utilities"
	"github.com/Sirupsen/logrus"
)

// GetRepositories returns a slice of repositories for the passed registry name
func GetRepositories(uri string) ([]string, error) {

	// Create and execute Get request for the catalog of repositores
	// https://github.com/docker/distribution/blob/master/docs/spec/api.md#catalog
	body, err := get(uri + "/_catalog")
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON into the catalog struct containing a slice of repositories
	c := catalog{}
	if err = json.Unmarshal(body, &c); err != nil {
		utils.Log.WithFields(logrus.Fields{
			"Error":         err,
			"Response Body": string(body),
		}).Error("Unable to unmarshal JSON!")
		return nil, err
	}

	return c.Repositories, err
}
