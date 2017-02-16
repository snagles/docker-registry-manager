package client

import (
	"encoding/json"
	"strings"

	"code.cloudfoundry.org/bytefmt"
	"github.com/DemonVex/docker-registry-manager/utilities"
	"github.com/Sirupsen/logrus"
)

// GetImage returns the image information for a given tag
// HEAD /v2/<name>/manifests/<reference>
func GetImage(uri string, repository string, tag string) (Image, error) {

	// Create and execute Get request
	body, err := get(uri + "/" + repository + "/manifests/" + tag)

	img := Image{}
	if err := json.Unmarshal(body, &img); err != nil {
		utils.Log.WithFields(logrus.Fields{
			"Error": err,
		}).Error("Unable to unmarshal JSON!")
		return img, err
	}

	// V1 compatibility is an escaped string, so convert it to JSON and then update the key
	for index, v1 := range img.History {
		v1JSON := V1Compatibility{}
		err = json.Unmarshal([]byte(v1.V1CompatibilityStr), &v1JSON)
		if err != nil {
			utils.Log.Error(err)
		}
		v1JSON.SizeStr = bytefmt.ByteSize(uint64(v1JSON.Size))

		// Update the image if we have any size information from the v1 compatibility
		if v1JSON.Size != 0 {
			img.ContainsV1Size = true
		}

		// Get first 8 characters for the short ID
		v1JSON.IDShort = v1JSON.ID[0:7]

		// Remove shell command
		v1JSON.ContainerConfig.CmdClean = strings.Replace(v1JSON.ContainerConfig.Cmd[0], "/bin/sh -c #(nop)", "", -1)
		img.History[index].V1Compatibility = v1JSON
	}

	// Update each FsLayer size
	for index, layer := range img.FsLayers {
		// Create and execute Get request
		response, err := head(uri + "/" + repository + "/blobs/" + layer.BlobSum)
		if err != nil {
			utils.Log.Error(err)
		} else {
			img.FsLayers[index].Size = response.ContentLength
			img.FsLayers[index].SizeStr = bytefmt.ByteSize(uint64(response.ContentLength))
		}
	}

	return img, nil
}
