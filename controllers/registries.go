package controllers

import (
	"github.com/DemonVex/docker-registry-manager/models/client"
	"github.com/DemonVex/docker-registry-manager/models/manager"
	"github.com/astaxie/beego"
)

// RegistriesController extends the beego.Controller type
type RegistriesController struct {
	beego.Controller
}

// Get returns the template for the registries page
func (c *RegistriesController) Get() {

	c.Data["registries"] = registry.Registries

	// Index template
	c.TplName = "registries.tpl"
}

func (c *RegistriesController) GetRegistryCount() {
	c.Data["registries"] = registry.Registries

	registryCount := struct {
		Count int
	}{
		len(registry.Registries),
	}
	c.Data["json"] = &registryCount
	c.ServeJSON()
}

// AddRegistry adds a registry to the active registry list from a form
func (c *RegistriesController) AddRegistry() {
	host := c.GetString("host")
	port := c.GetString("port")
	scheme := c.GetString("scheme")
	uri := scheme + "://" + host + ":" + port + "/v2"

	// Registry contains all identifying information for communicating with a registry
	err := registry.AddRegistry(uri)
	if err != nil {
		c.CustomAbort(404, err.Error())
	}
	c.Ctx.Redirect(302, "/registries")
}

// TestRegistryStatus responds with JSON containing the status of the registry
func (c *RegistriesController) TestRegistryStatus() {

	// Define the response
	var res struct {
		Error       string `json:"error, omitempty"`
		IsAvailable bool   `json:"is_available"`
	}

	host := c.GetString("host")
	port := c.GetString("port")
	scheme := c.GetString("scheme")
	uri := scheme + "://" + host + ":" + port + "/v2"

	// run the health check
	err := client.HealthCheck(uri)
	if err != nil {
		res.Error = err.Error()
		res.IsAvailable = false
		c.Data["json"] = &res
		c.ServeJSON()
	}

	res.Error = ""
	res.IsAvailable = true
	c.Data["json"] = &res
	c.ServeJSON()
}
