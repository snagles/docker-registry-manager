package dockerhub

import (
	manifestV2 "github.com/docker/distribution/manifest/schema2"
	"github.com/sirupsen/logrus"
	client "github.com/snagles/docker-registry-client/registry"
)

// URL for dockerhub connection
const URL = "https://registry.hub.docker.com"

// GetManifest retrieves the manifest for the given repo and tag name
func GetManifest(repoName string, tagName string) (*manifestV2.DeserializedManifest, error) {
	hub, err := client.New(URL, "", "")
	if err != nil {
		return nil, err
	}
	// add and retrieve oauth token for connections
	hub.Client.Transport = client.WrapTransport(hub.Client.Transport, URL, "", "")
	manifest, err := hub.Manifest("library/"+repoName, tagName)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error":           err.Error(),
			"Repository Name": repoName,
			"Tag Name":        tagName,
		}).Info("Failed to retrieve manifest information")
	}

	return manifest, err
}
