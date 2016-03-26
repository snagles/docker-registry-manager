package registry

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// TestGetRepositoriesByRegistry passes the ParseRegistry function a valid and invalid string
func TestGetRepositoriesByRegistry(t *testing.T) {

	// Create a new registry
	rURI := "https://host.domain:5000"
	r, err := ParseRegistry(rURI)

	Convey("When we pass an valid RegistryURI we should get back the registry type with no errors", t, func() {
		So(err, ShouldBeNil)
		So(r, ShouldNotBeEmpty)
	})

	err = GetRegistryStatus(rURI)
	// Test the response error
	Convey("When we get the registry status there should be no errors and the registry should be added to the map of active registries", t, func() {
		So(err, ShouldBeNil)
		So(ActiveRegistries, ShouldContainKey, "host.domain")
	})

	repos, err := GetRepositoriesByRegistry(r.Name)
	Convey("When we get the repositories for this registry there should be no errors and there should be a slice of repositories", t, func() {
		So(err, ShouldBeNil)
		So(repos, ShouldNotBeEmpty)
	})

}
