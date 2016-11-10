package client

import (
	"github.com/Sirupsen/logrus"
	"github.com/snagles/docker-registry-manager/utilities"
)

// HealthCheck takes in a registry URL and checks for communication errors
//
// Create and execute basic GET request to test if each registry can be reached
// To determine registry status we test the base registry route of /v2/ and check
// the HTTP response code for a 200 response (200 is a successful request)
func HealthCheck(uri string) error {

	// Create and execute a plain get request and check the http status code
	response, err := get(uri)
	if err != nil {
		// Notify of error
		utils.Log.WithFields(logrus.Fields{
			"Registry URLs": uri,
			"Error":         err,
			"HTTP Response": response,
			"Possible Fix":  "Check to see if your registry is up, and serving on the correct port with 'docker ps'.",
		}).Error("Get request to registry timed out/failed! Is the URL correct, and is the registry active?")
	}

	return err
}
