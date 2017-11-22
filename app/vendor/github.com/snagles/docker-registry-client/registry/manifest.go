package registry

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	manifestV1 "github.com/docker/distribution/manifest/schema1"
	manifestV2 "github.com/docker/distribution/manifest/schema2"
	digest "github.com/opencontainers/go-digest"
)

func (registry *Registry) Manifest(repository, reference string) (*manifestV2.DeserializedManifest, error) {
	url := registry.url("/v2/%s/manifests/%s", repository, reference)
	registry.Logf("registry.manifest.get url=%s repository=%s reference=%s", url, repository, reference)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", manifestV2.MediaTypeManifest)
	resp, err := registry.Client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	deserialized := &manifestV2.DeserializedManifest{}
	err = deserialized.UnmarshalJSON(body)
	return deserialized, err
}

func (registry *Registry) ManifestDigest(repository, reference string) (digest.Digest, error) {
	url := registry.url("/v2/%s/manifests/%s", repository, reference)
	registry.Logf("registry.manifest.head url=%s repository=%s reference=%s", url, repository, reference)

	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", manifestV2.MediaTypeManifest)
	resp, err := registry.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return digest.Parse(resp.Header.Get("Docker-Content-Digest"))
}

func (registry *Registry) HasManifest(repository, reference string) (bool, error) {
	checkURL := registry.url("/v2/%s/manifests/%s", repository, reference)
	registry.Logf("registry.manifest.head url=%s repository=%s reference=%s", checkURL, repository, reference)

	resp, err := registry.Client.Head(checkURL)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err == nil {
		return resp.StatusCode == http.StatusOK, nil
	}

	urlErr, ok := err.(*url.Error)
	if !ok {
		return false, err
	}
	httpErr, ok := urlErr.Err.(*HttpStatusError)
	if !ok {
		return false, err
	}

	if httpErr.Response.StatusCode == http.StatusNotFound {
		return false, nil
	}

	return false, err
}

func (registry *Registry) DeleteManifest(repository string, digest digest.Digest) error {
	url := registry.url("/v2/%s/manifests/%s", repository, digest)
	registry.Logf("registry.manifest.delete url=%s repository=%s reference=%s", url, repository, digest)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", manifestV2.MediaTypeManifest)
	resp, err := registry.Client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	return err
}

func (registry *Registry) PushManifest(repository, reference, mediaType string, payload []byte) (string, error) {
	url := registry.url("/v2/%s/manifests/%s", repository, reference)

	buffer := bytes.NewBuffer(payload)
	req, err := http.NewRequest("PUT", url, buffer)
	if err != nil {
		return "", err
	}
	req.Header.Set(http.CanonicalHeaderKey("Content-Type"), mediaType)
	resp, err := registry.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		digest := resp.Header.Get(http.CanonicalHeaderKey("Docker-Content-Digest"))
		return digest, nil
	}

	return "", fmt.Errorf("response status code is : %d", resp.StatusCode)

}

func (registry *Registry) PutManifest(repository, reference string, signedManifest *manifestV1.SignedManifest) error {
	url := registry.url("/v2/%s/manifests/%s", repository, reference)
	registry.Logf("registry.manifest.put url=%s repository=%s reference=%s", url, repository, reference)

	body, err := signedManifest.MarshalJSON()
	if err != nil {
		return err
	}

	buffer := bytes.NewBuffer(body)
	req, err := http.NewRequest("PUT", url, buffer)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", manifestV1.MediaTypeManifest)
	resp, err := registry.Client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	return err
}

func (registry *Registry) ManifestMetadata(repository string, digest digest.Digest) ([]byte, error) {
	url := registry.url("/v2/%s/blobs/%s", repository, digest.String())
	registry.Logf("registry.layer.check url=%s repository=%s digest=%s", url, repository, digest)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", manifestV1.MediaTypeManifest)
	resp, err := registry.Client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
