package utils

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// TestUpdateApp tests to see if all of the git functions execute successfully
func TestUpdateApp(t *testing.T) {

	_, err := UpdateApp()
	Convey("We should be able to get all of the local git ", t, func() {
		So(err, ShouldBeNil)
	})

}
