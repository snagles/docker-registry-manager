package manager

import (
	"testing"
	"time"

	"github.com/docker/distribution"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTagLastModifiedTime(t *testing.T) {
	setTime := time.Now().UTC()
	tag := Tag{
		V1Compatibility: new(V1Compatibility),
	}
	tag.V1Compatibility.History = []struct {
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
	Convey("Last modified time should be "+setTime.AddDate(0, 0, -1).String(), t, func(c C) {
		c.So(tag.LastModified(), ShouldResemble, setTime.AddDate(0, 0, -1))
	})

}
