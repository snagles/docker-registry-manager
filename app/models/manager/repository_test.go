package manager

import (
	"testing"
	"time"

	"github.com/docker/distribution"
	. "github.com/smartystreets/goconvey/convey"
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
	Convey("Last modified time should be "+setTime.AddDate(0, 0, -1).String(), t, func(c C) {
		c.So(repo.LastModified(), ShouldResemble, setTime.AddDate(0, 0, -1))
	})
}
