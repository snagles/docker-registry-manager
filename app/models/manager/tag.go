package manager

import (
	"time"

	manifestV2 "github.com/docker/distribution/manifest/schema2"
)

type Tag struct {
	*manifestV2.DeserializedManifest
	ID   string
	Name string
	*V1Compatibility
	Size int64
}

func (t *Tag) LastModified() (lastModified time.Time) {
	for _, history := range t.History {
		if history.Created.After(lastModified) {
			lastModified = history.Created
		}
	}
	return lastModified
}
