package manager

import (
	"fmt"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/snagles/docker-registry-manager/app/testutils"
)

func TestAddRegistry(t *testing.T) {
	baseurl, env := testutils.SetupRegistry(t)
	u, _ := url.Parse(baseurl)
	port, _ := strconv.Atoi(u.Port())
	r, err := AddRegistry(u.Scheme, u.Hostname(), "", "", port, 1*time.Minute, true)
	if err != nil {
		t.Fatalf("Failed to add test registry: %s", err)
	}

	if tr, ok := AllRegistries.Registries[fmt.Sprintf("%s:%v", r.Host, r.Port)]; ok {
		if tr.Status() != "UP" {
			t.Fatalf("Added registry status not up, reported as: %s", tr.Status())
		}
		env.Shutdown()
		if tr.Status() != "DOWN" {
			t.Fatalf("Added registry status not down, reported as: %s", tr.Status())
		}

	} else {
		t.Fatalf("Test registry not found in map of all registries: %s", r.URL)
	}
}
