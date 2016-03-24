package registry

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// TestParseRegistry passes the ParseRegistry function a valid and invalid string
func TestParseRegistry(t *testing.T) {

	validRegistryURI := "https://google.com:5000"
	// Create a registy type that contains the expected output from ParseRegistry
	expectedRegistryResponse := Registry{
		Name:   "google.com",
		Scheme: "https",
		Port:   "5000",
	}

	// Test the response error and use a deep equals comparison on the returned Registry
	registryResponse, err := ParseRegistry(validRegistryURI)
	Convey("When we pass a valid RegistryURI we should get back the registry type without an errors", t, func() {
		So(err, ShouldBeNil)
		So(registryResponse.Name, ShouldEqual, expectedRegistryResponse.Name)
		So(registryResponse.Scheme, ShouldEqual, expectedRegistryResponse.Scheme)
		So(registryResponse.Port, ShouldEqual, expectedRegistryResponse.Port)
	})

	invalidRegistryURI := "192.168.1.2:5000"
	registryResponse, err = ParseRegistry(invalidRegistryURI)
	// Test the response error
	Convey("When we pass an invalid RegistryURI we should get back the registry type with errors", t, func() {
		So(err, ShouldNotBeNil)
	})

}
