package manager

import (
	"testing"
	"time"

	"github.com/docker/distribution"
	"github.com/docker/distribution/manifest/schema2"
)

var setTime = time.Now().UTC()

func TestRepoLastModifiedTime(t *testing.T) {
	tag1 := Tag{
		V1Compatibility: new(V1Compatibility),
	}
	tag1.V1Compatibility.History = []struct {
		Created       time.Time                `json:"created"`
		Author        string                   `json:"author,omitempty"`
		CreatedBy     string                   `json:"created_by,omitempty"`
		Comment       string                   `json:"comment,omitempty"`
		EmptyLayer    bool                     `json:"empty_layer,omitempty"`
		ManifestLayer *distribution.Descriptor `json:"manifest_layer"`
		ShellType     string
		Commands      []Command
	}{
		{Created: setTime.AddDate(0, 0, -4)},
		{Created: setTime.AddDate(0, 0, -3)},
		{Created: setTime.AddDate(0, 0, -1)},
		{Created: setTime.AddDate(0, 0, -5)},
		{Created: setTime.AddDate(0, 0, -6)},
	}

	tag2 := Tag{
		V1Compatibility: new(V1Compatibility),
	}
	tag2.V1Compatibility.History = []struct {
		Created       time.Time                `json:"created"`
		Author        string                   `json:"author,omitempty"`
		CreatedBy     string                   `json:"created_by,omitempty"`
		Comment       string                   `json:"comment,omitempty"`
		EmptyLayer    bool                     `json:"empty_layer,omitempty"`
		ManifestLayer *distribution.Descriptor `json:"manifest_layer"`
		ShellType     string
		Commands      []Command
	}{
		{Created: setTime.AddDate(0, 0, -7)},
		{Created: setTime.AddDate(0, 0, -8)},
		{Created: setTime.AddDate(0, 0, -9)},
	}
	repo := Repository{
		Tags: map[string]*Tag{
			"tag1": &tag1,
			"tag2": &tag2,
		},
	}
	if repo.LastModified() != setTime.AddDate(0, 0, -1) {
		t.Errorf("Last modified time should be %s", setTime.AddDate(0, 0, -1))
	}
}

func TestRepoSize(t *testing.T) {
	tag1 := Tag{
		DeserializedManifest: new(schema2.DeserializedManifest),
	}
	tag1.Layers = []distribution.Descriptor{
		{Size: 500, Digest: "sha256:unique1"},
		{Size: 600, Digest: "sha256:unique2"},
		{Size: 700, Digest: "sha256:unique3"},
		{Size: 800, Digest: "sha256:unique4"},
		{Size: 900, Digest: "sha256:dup1"},
	}

	tag2 := Tag{
		DeserializedManifest: new(schema2.DeserializedManifest),
	}
	tag2.Layers = []distribution.Descriptor{
		{Size: 500, Digest: "sha256:unique5"},
		{Size: 600, Digest: "sha256:unique6"},
		{Size: 700, Digest: "sha256:unique7"},
		{Size: 800, Digest: "sha256:unique8"},
		{Size: 900, Digest: "sha256:dup1"},
	}
	repo := Repository{
		Tags: map[string]*Tag{
			"tag1": &tag1,
			"tag2": &tag2,
		},
	}

	if repo.Size() != 6100 {
		t.Errorf("Total repository size %v not equal to expected value: %v", repo.Size(), 6100)
	}
}
