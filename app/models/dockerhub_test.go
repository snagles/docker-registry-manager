package manager

import (
	"testing"

	test "github.com/snagles/docker-registry-manager/test/registry"
)

func TestHubGetManifest(t *testing.T) {
	for repo, tags := range test.Repositories {
		for _, tag := range tags {
			manifest, err := HubGetManifest(repo, tag)
			if err != nil || manifest == nil {
				if err != nil {
					t.Fatalf("Failed to retrieve %s manifest from docker hub: %s", tag, err)
				} else {
					t.Fatalf("Empty %s manifest retrieved from docker hub", tag)
				}
			}
		}
	}
}
