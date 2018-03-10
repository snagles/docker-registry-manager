package manager

import (
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/snagles/docker-registry-manager/app/testutils"
)

func TestNewRegistry(t *testing.T) {
	baseurl, env := testutils.SetupRegistry(t)
	u, _ := url.Parse(baseurl)
	port, _ := strconv.Atoi(u.Port())
	err := AllRegistries.AddRegistry(u.Scheme, u.Hostname(), "test", "", "", port, 1*time.Minute, true, true)
	if err != nil {
		t.Fatalf("Failed to add test registry: %s", err)
	}

	if tr, ok := AllRegistries.Registries["test"]; ok {
		if tr.Status() != "UP" {
			t.Fatalf("Added registry status not up, reported as: %s", tr.Status())
		}
		env.Shutdown()
	} else {
		t.Fatalf("Test registry not found in map of all registries: %s", "test")
	}
}
