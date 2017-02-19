package client

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/snagles/docker-registry-manager/utilities"
)

// GetTags returns a slice of tags for a given repository and registry
// https://github.com/docker/distribution/blob/master/docs/spec/api.md#listing-image-tags
func GetTags(uri string, repository string) ([]string, error) {

	// Create and execute Get request
	body, err := get(uri + "/" + repository + "/tags/list")
	if err != nil {
		return nil, err
	}

	l := list{}
	// Unmarshal JSON into the tag response struct containing an array of tags
	if err := json.Unmarshal(body, &l); err != nil {
		utils.Log.WithFields(logrus.Fields{
			"Error":         err,
			"Response Body": string(body),
		}).Error("Unable to unmarshal JSON!")
		return l.Tags, err
	}
	return l.Tags, nil
}

// DeleteTag deletes the tag by first getting its docker-content-digest, and then using
// the digest received the function deletes the manifest
//
// Documentation:
// DELETE	/v2/<name>/manifests/<reference>	Manifest	Delete the manifest identified by name and reference. Note that a manifest can only be deleted by digest.
func DeleteTag(uri string, repository string, tag string) (bool, error) {

	// Check if the tag exists. If it does not we cannot get the digest from it
	client := &http.Client{}
	req, err := http.NewRequest("HEAD", uri+"/"+repository+"/manifests/"+tag, nil)
	if err != nil {
		utils.Log.Debug(err)
	}

	// Note When deleting a manifest from a registry version 2.3 or later, the following header must be used when HEAD or GET-ing the manifest to obtain the correct digest to delete:
	// Accept: application/vnd.docker.distribution.manifest.v2+json
	req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		utils.Log.WithFields(logrus.Fields{
			"Error":    err,
			"Tag":      tag,
			"Response": resp,
		}).Error("Could not delete tag! Could not head the tag.")
		return false, err
	}

	// Make sure the digest exists in the header. If it does, attempt the deletion
	if _, ok := resp.Header["Docker-Content-Digest"]; ok {

		if len(resp.Header["Docker-Content-Digest"]) > 0 {
			// Create and execute DELETE request
			digest := resp.Header["Docker-Content-Digest"][0]
			client := &http.Client{}
			req, _ := http.NewRequest("DELETE", uri+"/"+repository+"/manifests/"+digest, nil)
			req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")
			resp, err = client.Do(req)
			if err != nil || (resp.StatusCode < 200 || resp.StatusCode >= 300) {
				body, readErr := ioutil.ReadAll(resp.Body)
				utils.Log.WithFields(logrus.Fields{
					"Connection Error": err,
					"Error Response":   string(body),
					"Tag":              tag,
					"Read Error":       readErr,
				}).Error("Could not delete tag! Received: " + string(body))
				err = errors.New(string(body))

				return false, err
			}
			return true, nil
		}
	}
	// Error if there was nothing in the Docker-Content-Digest field
	utils.Log.WithFields(logrus.Fields{
		"Error": errors.New("No digest gotten from response header"),
		"Tag":   tag,
	}).Error("Could not delete tag!")
	return false, errors.New("No digest gotten from response header")

}
