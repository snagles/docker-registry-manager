package manager

import (
	"github.com/Sirupsen/logrus"
	manifestV2 "github.com/docker/distribution/manifest/schema2"
	client "github.com/heroku/docker-registry-client/registry"
)

const (
	//HubURL contains the registry url for dockerhub
	HubURL = "https://registry.hub.docker.com"
)

// HubUser set to empty so basic auth is not used
var HubUser = "" // anonymous
// HubPassword set to empty so basic auth is not used
var HubPassword = "" // anonymous

// HubGetManifest retrieves the manifest for the given repo and tag name
func HubGetManifest(repoName string, tagName string) (*manifestV2.DeserializedManifest, error) {
	hub, err := client.New(HubURL, HubUser, HubPassword)
	if err != nil {
		return nil, err
	}
	// get oauth token and use it as transport
	hub.Client.Transport = client.WrapTransport(hub.Client.Transport, HubURL, HubUser, HubPassword)
	manifest, err := hub.ManifestV2("library/"+repoName, tagName)
	// Using v2 required getting the manifest then retrieving the blob
	// for the config digest
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error":           err.Error(),
			"Repository Name": repoName,
			"Tag Name":        tagName,
		}).Error("Failed to retrieve manifest information")
	}

	return manifest, err
}
