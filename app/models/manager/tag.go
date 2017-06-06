package manager

import (
	"time"

	manifestV1 "github.com/docker/distribution/manifest/schema1"
	manifestV2 "github.com/docker/distribution/manifest/schema2"
	"github.com/snagles/docker-registry-manager/utils"
)

type Tag struct {
	V1        *manifestV1.SignedManifest
	V2        *manifestV2.DeserializedManifest
	ID        string
	Name      string
	Histories []V1Compatibility
}

func (t *Tag) LastModified() time.Time {
	var lastModified time.Time
	for _, history := range t.Histories {
		if history.Created.After(lastModified) {
			lastModified = history.Created
		}
	}
	return lastModified
}

func (t *Tag) LastModifiedTimeAgo() string {
	lastModified := t.LastModified()
	return utils.TimeAgo(lastModified)
}

func (t *Tag) LayerCount() int {
	return len(t.V1.FSLayers)
}

func (t *Tag) HistoriesOrdered() []V1Compatibility {
	ordered := []V1Compatibility{}
	for i := len(t.Histories) - 1; i >= 0; i-- {
		ordered = append(ordered, t.Histories[i])
	}
	return ordered
}
