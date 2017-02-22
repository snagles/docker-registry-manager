package client

import (
	"encoding/json"
	"regexp"

	"github.com/DemonVex/docker-registry-manager/utilities"
	"github.com/Sirupsen/logrus"
)

// GetRepositories returns a slice of repositories for the passed registry name
func GetRepositories(uri string) ([]string, error) {
	repositories := []string{}
	path := "/_catalog"

	r, err := regexp.Compile("</v2(.*)>;")
	if err != nil {
		return nil, err
	}

	// Default limit is 100 so we need fetch data more than once
	for {
		// Create and execute Get request for the catalog of repositores
		// https://github.com/docker/distribution/blob/master/docs/spec/api.md#catalog
		body, headers, err := getWithHeaders(uri + path)

		// Unmarshal JSON into the catalog struct containing a slice of repositories
		c := catalog{}
		if err = json.Unmarshal(body, &c); err != nil {
			utils.Log.WithFields(logrus.Fields{
				"Error":         err,
				"Response Body": string(body),
			}).Error("Unable to unmarshal JSON!")
			return nil, err
		}
		repositories = append(repositories, c.Repositories...)

		//If there are more repositories than limit we will receive a header Link
		// Link: </v2/_catalog?n=<last n value>&last=<last entry from response>>; rel="next"
		link := r.FindStringSubmatch(headers.Get("Link"))
		if len(link) == 0 {
			break
		}
		path = link[1]
	}

	return repositories, nil
}
