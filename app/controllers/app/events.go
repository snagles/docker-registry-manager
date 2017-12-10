package app

import (
	"encoding/json"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/snagles/docker-registry-manager/app/models"
)

// EventsController handles all posting of envelopes from the registries configured endpoints
// in addition to all retrieval of the in memory store of posted registry events
type EventsController struct {
	beego.Controller
}

func (c *EventsController) Get() {
	manager.AllEvents.Lock()
	defer manager.AllEvents.Unlock()
	c.Data["events"] = manager.AllEvents.Events
	c.TplName = "events.tpl"
}

// GetRegistryEvents returns the JSON representation of all events for the given registry
func (c *EventsController) GetRegistryEvents() {
	// Get the registry
	registryName := c.Ctx.Input.Param(":registryName")

	manager.AllEvents.Lock()
	defer manager.AllEvents.Unlock()

	// make sure the registry exists in the map
	if _, ok := manager.AllEvents.Events[registryName]; !ok {
		c.CustomAbort(404, "No events found for registry.")
		return
	}
	c.Data["json"] = manager.AllEvents.Events[registryName]
	c.ServeJSON()
}

// GetRegistryEventID uses the passed registry and unique event id to return the json
// representation of the event
func (c *EventsController) GetRegistryEventID() {
	// Get the registry
	registryName := c.Ctx.Input.Param(":registryName")

	manager.AllEvents.Lock()
	defer manager.AllEvents.Unlock()

	// make sure the registry exists in the map
	if _, ok := manager.AllEvents.Events[registryName]; !ok {
		c.CustomAbort(404, "No events found for registry.")
		return
	}

	// get  the registry id
	eventID := c.Ctx.Input.Param(":eventID")
	if _, ok := manager.AllEvents.Events[registryName][eventID]; !ok {
		c.CustomAbort(404, "Event ID not found.")
		return
	}
	c.Data["json"] = manager.AllEvents.Events[registryName][eventID]
	c.ServeJSON()
}

// GetEvents returns the JSON representation of all events received from ALL of the registries
func (c *EventsController) GetEvents() {
	manager.AllEvents.Lock()
	defer manager.AllEvents.Unlock()
	c.Data["json"] = manager.AllEvents.Events
	c.ServeJSON()
}

// PostEnvelope is th endpoint used by the notification systems of the given registry.
// see https://docs.docker.com/registry/notifications/#configuration for details
func (c *EventsController) PostEnvelope() {
	e := manager.Envelope{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &e); err != nil {
		c.CustomAbort(http.StatusNotFound, "Unable to parse envelope.")
	}
	e.Process()
	c.CustomAbort(200, "")
}
