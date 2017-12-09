package manager

import (
	"time"

	manifestV2 "github.com/docker/distribution/manifest/schema2"
)

// Tag contains all v1 compatibility and manifest information return by the registry
type Tag struct {
	*manifestV2.DeserializedManifest
	ID   string
	Name string
	*V1Compatibility
	Size int64
}

// LastModified returns the latest last modified time using the history fields
func (t *Tag) LastModified() (lastModified time.Time) {
	for _, history := range t.History {
		if history.Created.After(lastModified) {
			lastModified = history.Created
		}
	}
	return lastModified
}
