package manager

import "sync"

// AllEvents contains all events received for each individual registry, to read
// or write from the list of events the mutex needs to be used to ensure consistency
var AllEvents Events

func init() {
	AllEvents = Events{
		Events: make(map[string]map[string]Event),
	}
}

// Events contains event information in the structure registry->event_id->event
type Events struct {
	sync.Mutex
	Events map[string]map[string]Event
}

// Envelope contains all grouped events sent by the registries.
// see https://docs.docker.com/registry/notifications/#envelope for more details
type Envelope struct {
	Events []Event
}

// Process adds each individual event to the map of all events
func (e *Envelope) Process() {
	AllEvents.Lock()
	for _, event := range e.Events {
		if event.Action == "push" && event.Request.Useragent != "Go-http-client/1.1" && event.Request.Method != "HEAD" {
			if _, ok := AllEvents.Events[event.Request.Host]; !ok {
				AllEvents.Events[event.Request.Host] = make(map[string]Event)
			}
			// Add the event
			AllEvents.Events[event.Request.Host][event.ID] = event
		}
	}
	AllEvents.Unlock()
}

// Event contains all information for the given event (e.g pull/pushes to the registry)
// see https://docs.docker.com/registry/notifications/#events for more details
type Event struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Action    string `json:"action"`
	Target    struct {
		MediaType  string `json:"mediaType"`
		Size       int    `json:"size"`
		Digest     string `json:"digest"`
		Length     int    `json:"length"`
		Repository string `json:"repository"`
		URL        string `json:"url"`
		Tag        string `json:"tag"`
	} `json:"target"`
	Request struct {
		ID        string `json:"id"`
		Addr      string `json:"addr"`
		Host      string `json:"host"`
		Method    string `json:"method"`
		Useragent string `json:"useragent"`
	} `json:"request"`
	Actor struct {
	} `json:"actor"`
	Source struct {
		Addr       string `json:"addr"`
		InstanceID string `json:"instanceID"`
	} `json:"source"`
}
