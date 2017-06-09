package manager

import (
	"time"

	"github.com/snagles/docker-registry-manager/utils"
)

type Repository struct {
	Name string
	Tags map[string]*Tag
}

func (r *Repository) LastModified() time.Time {
	var lastModified time.Time
	for _, tag := range r.Tags {
		for _, history := range tag.History {
			if history.Created.After(lastModified) {
				lastModified = history.Created
			}
		}
	}
	return lastModified
}

func (r *Repository) LastModifiedTimeAgo() string {
	lastModified := r.LastModified()
	return utils.TimeAgo(lastModified)
}
