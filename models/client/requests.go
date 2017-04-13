package client

import (
	"io/ioutil"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/snagles/docker-registry-manager/utilities"
)

// Helper function for get requests
func get(uri string) ([]byte, error) {
	body, _, err := getWithHeaders(uri)
	return body, err
}

func getWithHeaders(uri string) ([]byte, http.Header, error) {
	response, err := http.Get(uri)
	if err != nil {
		logRequestError(uri, err)
		return nil, nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil || response.StatusCode != 200 {
		utils.Log.WithFields(logrus.Fields{
			"Error": err,
			"Body":  body,
		}).Error("Unable to read response body returned from the registry!")
		return nil, nil, err
	}
	return body, response.Header, err
}

// head function for head requests
func head(uri string) (*http.Response, error) {
	response, err := http.Head(uri)
	if err != nil {
		logRequestError(uri, err)
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		utils.Log.WithFields(logrus.Fields{
			"Registry URL": uri,
			"Status Code":  response.StatusCode,
			"Response":     response,
			"Possible Fix": "Check to see if your registry is up, and serving on the correct port with 'docker ps'. ",
		}).Error("Get request to registry failed for the endpoint! Is your registry active?")
		return nil, err
	}
	return response, err
}

func logRequestError(uri string, err error) {
	utils.Log.WithFields(logrus.Fields{
		"Registry URL": uri,
		"Error":        err,
		"Possible Fix": "Check to see if your registry is up, and serving on the correct port with 'docker ps'. ",
	}).Error("Get request to registry failed for the endpoint! Is your registry active?")
}
