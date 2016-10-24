package controllers

import (
	"github.com/astaxie/beego"
	"github.com/snagles/docker-registry-manager/models/registry"
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
	registry.AddRegistry(uri)
	c.Ctx.Redirect(302, "/registries")
}
