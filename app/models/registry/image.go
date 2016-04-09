package registry

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
)

// Image contains all information related to the image
type Image struct {
	Name          string
	Tag           string
	SchemaVersion int
	Architecture  string
	TagID         uint
	History       []History `json:"history"`
	FsLayers      []struct {
		BlobSum string `json:"blobSum"`
	} `json:"fsLayers"`
}

// History contains the v1 compatibility string and marshaled json
type History struct {
	V1CompatibilityStr string `json:"V1Compatibility"`
	V1Compatibility    V1Compatibility
}

// V1Compatibility contains all information grabbed from the V1Compatibility field from registry v1
type V1Compatibility struct {
	ID              string    `json:"id"`
	Parent          string    `json:"parent"`
	Created         time.Time `json:"created"`
	Container       string    `json:"container"`
	ContainerConfig struct {
		Hostname        string        `json:"Hostname"`
		Domainname      string        `json:"Domainname"`
		User            string        `json:"User"`
		AttachStdin     bool          `json:"AttachStdin"`
		AttachStdout    bool          `json:"AttachStdout"`
		AttachStderr    bool          `json:"AttachStderr"`
		ExposedPorts    interface{}   `json:"ExposedPorts"`
		PublishService  string        `json:"PublishService"`
		Tty             bool          `json:"Tty"`
		OpenStdin       bool          `json:"OpenStdin"`
		StdinOnce       bool          `json:"StdinOnce"`
		Env             []string      `json:"Env"`
		Cmd             []string      `json:"Cmd"`
		Image           string        `json:"Image"`
		Volumes         interface{}   `json:"Volumes"`
		VolumeDriver    string        `json:"VolumeDriver"`
		WorkingDir      string        `json:"WorkingDir"`
		Entrypoint      interface{}   `json:"Entrypoint"`
		NetworkDisabled bool          `json:"NetworkDisabled"`
		MacAddress      string        `json:"MacAddress"`
		OnBuild         []interface{} `json:"OnBuild"`
		Labels          struct {
		} `json:"Labels"`
	} `json:"container_config"`
	DockerVersion string `json:"docker_version"`
	Config        struct {
		Hostname        string        `json:"Hostname"`
		Domainname      string        `json:"Domainname"`
		User            string        `json:"User"`
		AttachStdin     bool          `json:"AttachStdin"`
		AttachStdout    bool          `json:"AttachStdout"`
		AttachStderr    bool          `json:"AttachStderr"`
		ExposedPorts    interface{}   `json:"ExposedPorts"`
		PublishService  string        `json:"PublishService"`
		Tty             bool          `json:"Tty"`
		OpenStdin       bool          `json:"OpenStdin"`
		StdinOnce       bool          `json:"StdinOnce"`
		Env             []string      `json:"Env"`
		Cmd             []string      `json:"Cmd"`
		Image           string        `json:"Image"`
		Volumes         interface{}   `json:"Volumes"`
		VolumeDriver    string        `json:"VolumeDriver"`
		WorkingDir      string        `json:"WorkingDir"`
		Entrypoint      interface{}   `json:"Entrypoint"`
		NetworkDisabled bool          `json:"NetworkDisabled"`
		MacAddress      string        `json:"MacAddress"`
		OnBuild         []interface{} `json:"OnBuild"`
		Labels          struct {
		} `json:"Labels"`
	} `json:"config"`
	Architecture string `json:"architecture"`
	Os           string `json:"os"`
	Size         int    `json:"Size"`
}

// GetImage returns the image information for a given tag
// HEAD /v2/<name>/manifests/<reference>
/*
	"name": <name>,
	"tag": <tag>,
	"fsLayers": [
		 {
				"blobSum": "<digest>"
		 },
		 ...
	 ]
	],
	"history": <v1 images>,
	"signature": <JWS>
*/
func GetImage(registryName string, repositoryName string, tagName string) (Image, error) {

	// Check if the registry is listed as active
	if _, ok := ActiveRegistries[registryName]; !ok {
		return Image{}, errors.New(registryName + " was not found within the active list of registries.")
	}
	r := ActiveRegistries[registryName]

	// Create and execute Get request
	response, err := http.Get(r.GetURI() + "/" + repositoryName + "/manifests/" + tagName)

	if response.StatusCode != 200 {
		log.WithFields(log.Fields{
			"Error":       err,
			"Status Code": response.StatusCode,
			"Response":    response,
		}).Error("Did not receive an ok status code!")
		return Image{}, err
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
		return Image{}, err
	}
	img := Image{}
	if err := json.Unmarshal(body, &img); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Unable to unmarshal JSON!")
		return Image{}, err
	}

	// V1 compatibility is an escape string, so convert it to JSON and then update the key
	for index, v1 := range img.History {
		v1JSON := V1Compatibility{}
		err = json.Unmarshal([]byte(v1.V1CompatibilityStr), &v1JSON)
		if err != nil {
			log.Error(err)
		}
		img.History[index].V1Compatibility = v1JSON
	}

	return img, nil
}
