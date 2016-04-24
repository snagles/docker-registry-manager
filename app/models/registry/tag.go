package registry

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/pivotal-golang/bytefmt"
	"github.com/stefannaglee/docker-registry-manager/app/utilities"
)

// Tags contains a slice of tags for the given repository
// https://github.com/docker/distribution/blob/master/docs/spec/api.md#listing-image-tags
type Tags struct {
	Name string
	Tags []string
}

type TagForView struct {
	ID              string
	Name            string
	UpdatedTime     time.Time
	UpdatedTimeUnix int64
	TimeAgo         string
	Layers          int
	Size            string
}

// TagsForView contains a slice of TagsForView with the methods required to sort
type TagsForView []TagForView

func (slice TagsForView) Len() int {
	return len(slice)
}

func (slice TagsForView) Less(i, j int) bool {
	return slice[i].Name < slice[j].Name
}

func (slice TagsForView) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// GetTagsForView returns the sanitized tag structs with the required information for the tags template
func GetTagsForView(registryName string, repositoryName string) (TagsForView, error) {
	tagObj, err := GetTags(registryName, repositoryName)
	tags := TagsForView{}

	tagChan := make(chan TagForView)

	// Loop through each tag to build the TagForView type
	for _, tagName := range tagObj.Tags {

		go func(tagName string) {
			// Created a new tag for view type to fill
			t := TagForView{}
			var tempSize int64
			var maxTime time.Time

			// Get the image information for each tag
			img, _ := GetImage(registryName, repositoryName, tagName)

			for _, layer := range img.FsLayers {

				// Check if the registry is listed as active
				r := ActiveRegistries[registryName]
				// Create and execute Get request
				response, _ := http.Head(r.GetURI() + "/" + repositoryName + "/blobs/" + layer.BlobSum)
				if err != nil {
					utils.Log.Error(err)
				}
				tempSize += response.ContentLength
			}
			// Get the latest creation time and total the size for the tag image
			for _, history := range img.History {
				if history.V1Compatibility.Created.After(maxTime) {
					maxTime = history.V1Compatibility.Created
				}
			}

			// Set the fields
			t.Size = bytefmt.ByteSize(uint64(tempSize))
			t.UpdatedTime = maxTime
			t.UpdatedTimeUnix = maxTime.Unix()
			t.Layers = len(img.History)
			t.Name = tagName
			t.TimeAgo = utils.TimeAgo(maxTime)

			// Append to the tags list that will be passed to the template
			tagChan <- t

		}(tagName)

	}

	var TagInformation TagsForView
	// Wait for each of the requests and append to the returned tag information
	for i := 0; i < len(tagObj.Tags); i++ {
		tag := <-tagChan
		TagInformation = append(TagInformation, tag)
	}
	close(tagChan)
	sort.Sort(sort.Reverse(tags))

	return TagInformation, err
}

// GetTags returns a slice of tags for a given repository and registry
func GetTags(registryName string, repositoryName string) (Tags, error) {

	repositoryName, _ = url.QueryUnescape(repositoryName)

	// Check if the registry is listed as active
	if _, ok := ActiveRegistries[registryName]; !ok {
		return Tags{}, errors.New(registryName + " was not found within the active list of registries.")
	}
	r := ActiveRegistries[registryName]

	// Create and execute Get request
	response, err := http.Get(r.GetURI() + "/" + repositoryName + "/tags/list")
	if err != nil {
		utils.Log.WithFields(logrus.Fields{
			"Registry URL": string(r.GetURI()),
			"Error":        err,
			"Possible Fix": "Check to see if your registry is up, and serving on the correct port with 'docker ps'. ",
		}).Error("Get request to registry failed for the tags endpoint.")
		return Tags{}, err
	}

	// Check Status code
	if response.StatusCode != 200 {
		utils.Log.WithFields(logrus.Fields{
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
		utils.Log.WithFields(logrus.Fields{
			"Error": err,
			"Body":  body,
		}).Error("Unable to read response into body!")
		return Tags{}, err
	}
	ts := Tags{}
	// Unmarshal JSON into the tag response struct containing an array of tags
	if err := json.Unmarshal(body, &ts); err != nil {
		utils.Log.WithFields(logrus.Fields{
			"Error":         err,
			"Response Body": string(body),
		}).Error("Unable to unmarshal JSON!")
		return ts, err
	}
	return ts, nil
}

// DeleteTag deletes the tag by first getting its docker-content-digest, and then using
// the digest received the function deletes the manifest
//
// Documentation:
// DELETE	/v2/<name>/manifests/<reference>	Manifest	Delete the manifest identified by name and reference. Note that a manifest can only be deleted by digest.
func DeleteTag(registryName string, repositoryName string, tag string) (bool, error) {

	repositoryName, _ = url.QueryUnescape(repositoryName)

	// Check if the registry is listed as active
	if _, ok := ActiveRegistries[registryName]; !ok {
		return false, errors.New(registryName + " was not found within the active list of registries.")
	}
	r := ActiveRegistries[registryName]

	// Check if the tag exists. If it does not we cannot get the digest from it
	client := &http.Client{}
	req, _ := http.NewRequest("HEAD", r.GetURI()+"/"+repositoryName+"/manifests/"+tag, nil)

	// Note When deleting a manifest from a registry version 2.3 or later, the following header must be used when HEAD or GET-ing the manifest to obtain the correct digest to delete:
	// Accept: application/vnd.docker.distribution.manifest.v2+json
	req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")

	// Execute the request
	resp, existsErr := client.Do(req)
	if existsErr != nil {
		utils.Log.WithFields(logrus.Fields{
			"Request":  resp.Request,
			"Error":    existsErr,
			"Tag":      tag,
			"Response": resp,
		}).Error("Could not delete tag! Could not head the tag.")
		return false, existsErr
	}

	// Make sure the digest exists in the header. If it does, attempt the deletion
	if _, ok := resp.Header["Docker-Content-Digest"]; ok {

		if len(resp.Header["Docker-Content-Digest"]) > 0 {
			// Create and execute DELETE request
			digest := resp.Header["Docker-Content-Digest"][0]
			client := &http.Client{}
			req, _ := http.NewRequest("DELETE", r.GetURI()+"/"+repositoryName+"/manifests/"+digest, nil)
			req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")
			resp, err := client.Do(req)
			if err != nil || resp.StatusCode != 200 {
				utils.Log.WithFields(logrus.Fields{
					"Error":    err,
					"Tag":      tag,
					"Response": resp,
				}).Error("Could not delete tag!")
				return false, err
			}
		}

		// Error if there was nothing in the Docker-Content-Digest field
		utils.Log.WithFields(logrus.Fields{
			"Error":    errors.New("No digest gotten from response header"),
			"Tag":      tag,
			"Response": resp,
		}).Error("Could not delete tag!")
		return false, errors.New("No digest gotten from response header")
	}

	return true, nil
}

// Get tag returns a TagForView based on the passed tag name
func GetTag(registryName string, repositoryName string, tagName string) (TagForView, error) {

	// Created a new tag for view type to fill
	t := TagForView{}
	var tempSize int64
	var maxTime time.Time

	// Get the image information for each tag
	img, _ := GetImage(registryName, repositoryName, tagName)

	for _, layer := range img.FsLayers {

		// Check if the registry is listed as active
		r := ActiveRegistries[registryName]
		// Create and execute Get request
		response, err := http.Head(r.GetURI() + "/" + repositoryName + "/blobs/" + layer.BlobSum)
		if err != nil {
			utils.Log.Error(err)
			return t, err
		}
		tempSize += response.ContentLength
	}
	// Get the latest creation time and total the size for the tag image
	for _, history := range img.History {
		if history.V1Compatibility.Created.After(maxTime) {
			maxTime = history.V1Compatibility.Created
		}
	}

	// Set the fields
	t.Size = bytefmt.ByteSize(uint64(tempSize))
	t.UpdatedTime = maxTime
	t.UpdatedTimeUnix = maxTime.Unix()
	t.Layers = len(img.History)
	t.Name = tagName
	t.TimeAgo = utils.TimeAgo(maxTime)

	return t, nil

}
