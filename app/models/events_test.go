package manager

import (
	"encoding/json"
	"testing"
)

func TestEvents(t *testing.T) {
	e := Envelope{}
	err := json.Unmarshal([]byte(envelope), &e)
	if err != nil {
		t.Errorf("Failed to unmarshal test envelope into Envelope struct: %s", err.Error())
	}

	e.Process()

	if _, ok := AllEvents.Events["192.168.100.227:5000"]; !ok {
		t.Errorf("AllEvents does not contain the correct request host name: %s", "192.168.100.227:5000")
	}

	if _, ok := AllEvents.Events["192.168.100.227:5000"]["320678d8-ca14-430f-8bb6-4ca139cd83f7"]; !ok {
		t.Errorf("AllEvents does not contain the correct unique identifier %s", "320678d8-ca14-430f-8bb6-4ca139cd83f7")
	}

	if AllEvents.Events["192.168.100.227:5000"]["320678d8-ca14-430f-8bb6-4ca139cd83f7"].Action != "push" {
		t.Error("AllEvents does not contain the correct push action")
	}
}

var envelope = `{
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
      }
   ]
}
`
