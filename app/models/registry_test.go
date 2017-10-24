package manager

import (
	"encoding/json"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/snagles/docker-registry-manager/test/registrytest"
)

func TestAddRegistry(t *testing.T) {
	tr := registrytest.New("https.yml")
	registrytest.Start(tr)
	time.Sleep(2 * time.Second)

	// Test valid form
	r, err := AddRegistry("http", "localhost", "", "", 5000, 10*time.Second, true)
	if err != nil {
		t.Error("Failed to add registry: " + err.Error())
	}
	if r == nil {
		t.Error("Returned registry does not contain expected values")
	}
	registrytest.Stop(tr)
}

func TestTagCount(t *testing.T) {
	registry := createTestRegistry()
	if registry.TagCount() != 1 {
		t.Error("Expected 1 tag, received: " + strconv.Itoa(registry.TagCount()))
	}
}

func TestIP(t *testing.T) {
	registry := createTestRegistry()
	registry.Host, _ = os.Hostname()
	if registry.IP() == "" {
		t.Error("Failed to get ip for hostname")
	}
}

func TestPushes(t *testing.T) {
	e := Envelope{}
	err := json.Unmarshal([]byte(testPushesPulls), &e)
	if err != nil {
		logrus.Fatalf("Failed to unmarshal test envelope into Envelope struct: %s", err.Error())
	}

	e.Process()
	r := createTestRegistry()
	if r.Pushes() != 1 {
		t.Errorf("Invalid number of pushes returned for %v, expected 1", r.Host)
	}
}

func TestPulls(t *testing.T) {
	e := Envelope{}
	err := json.Unmarshal([]byte(testPushesPulls), &e)
	if err != nil {
		logrus.Fatalf("Failed to unmarshal test envelope into Envelope struct: %s", err.Error())
	}

	e.Process()

	r := createTestRegistry()
	if r.Pulls() != 1 {
		t.Errorf("Invalid number of pulls returned for %v, expected 1", r.Host)
	}
}

func createTestRegistry() Registry {
	registry := Registry{
		Name:    "192.168.100.227:5000",
		Host:    "192.168.100.227:5000",
		Scheme:  "https",
		Version: "v2",
		Port:    5000,
	}
	repository := Repository{
		Name: "testrepo",
	}

	repository.Tags = make(map[string]*Tag)
	tag := Tag{
		ID:   "1",
		Name: "testTag",
		Size: 400,
	}
	repository.Tags["testTag"] = &tag
	registry.Repositories = make(map[string]*Repository)
	registry.Repositories["testRepo"] = &repository

	return registry
}

var testPushesPulls = `{
   "events": [
      {
         "id": "320678d8-ca14-430f-8bb6-4ca139cd83f7",
         "timestamp": "2016-03-09T14:44:26.402973972-08:00",
         "action": "push",
         "target": {
            "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
            "size": 708,
            "digest": "sha256:fea8895f450959fa676bcc1df0611ea93823a735a01205fd8622846041d0c7cf",
            "length": 708,
            "repository": "hello-world",
            "url": "http://192.168.100.227:5000/v2/hello-world/manifests/sha256:fea8895f450959fa676bcc1df0611ea93823a735a01205fd8622846041d0c7cf",
            "tag": "latest"
         },
         "request": {
            "id": "6df24a34-0959-4923-81ca-14f09767db19",
            "addr": "192.168.64.11:42961",
            "host": "192.168.100.227:5000",
            "method": "GET",
            "useragent": "curl/7.38.0"
         },
         "actor": {},
         "source": {
            "addr": "xtal.local:5000",
            "instanceID": "a53db899-3b4b-4a62-a067-8dd013beaca4"
         }
      },
			{
				 "id": "320678d8-ca14-430f-8bb6-4ca139cd8300",
				 "timestamp": "2016-03-09T14:44:26.402973972-08:00",
				 "action": "pull",
				 "target": {
						"mediaType": "application/vnd.docker.distribution.manifest.v2+json",
						"size": 708,
						"digest": "sha256:fea8895f450959fa676bcc1df0611ea93823a735a01205fd8622846041d0c7cf",
						"length": 708,
						"repository": "hello-world",
						"url": "http://192.168.100.227:5000/v2/hello-world/manifests/sha256:fea8895f450959fa676bcc1df0611ea93823a735a01205fd8622846041d0c7cf",
						"tag": "latest"
				 },
				 "request": {
						"id": "6df24a34-0959-4923-81ca-14f09767db00",
						"addr": "192.168.64.11:42961",
						"host": "192.168.100.227:5000",
						"method": "GET",
						"useragent": "curl/7.38.0"
				 },
				 "actor": {},
				 "source": {
						"addr": "xtal.local:5000",
						"instanceID": "a53db899-3b4b-4a62-a067-8dd013beac00"
				 }
			}
   ]
}
`
