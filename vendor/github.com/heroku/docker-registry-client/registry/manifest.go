package registry

import (
	"io/ioutil"
	"net/http"

	manifestV1 "github.com/docker/distribution/manifest/schema1"
	manifestV2 "github.com/docker/distribution/manifest/schema2"
	digest "github.com/opencontainers/go-digest"
)

func (registry *Registry) ManifestV2(repository, reference string) (*manifestV2.DeserializedManifest, error) {
	url := registry.url("/v2/%s/manifests/%s", repository, reference)
	registry.Logf("registry.manifest.get url=%s repository=%s reference=%s", url, repository, reference)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", manifestV2.MediaTypeManifest)
	resp, err := registry.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	deserialized := &manifestV2.DeserializedManifest{}
	err = deserialized.UnmarshalJSON(body)
	if err != nil {
		return nil, err
	}
	return deserialized, nil
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
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return "", err
	}
	return digest.Digest(resp.Header.Get("Docker-Content-Digest")), nil
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
	if err != nil {
		return err
	}
	return nil
}

func (registry *Registry) PutManifest(repository, reference string) error {
	url := registry.url("/v2/%s/manifests/%s", repository, reference)
	registry.Logf("registry.manifest.put url=%s repository=%s reference=%s", url, repository, reference)

	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", manifestV2.MediaTypeManifest)
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
