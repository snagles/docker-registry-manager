package controllers

import (
	"github.com/astaxie/beego"
	"github.com/stefannaglee/docker-registry-manager/app/models/registry"
	"github.com/stefannaglee/docker-registry-manager/app/utilities"
)

// RegistriesController extends the beego.Controller type
type RegistriesController struct {
	beego.Controller
}

// Get returns the template for the registries page
func (c *RegistriesController) Get() {

	c.Data["registries"] = registry.ActiveRegistries

	// Index template
	c.TplName = "registries.tpl"
}

func (c *RegistriesController) GetRegistryCount() {
	c.Data["registries"] = registry.ActiveRegistries

	registryCount := struct {
		Count int
	}{
		len(registry.ActiveRegistries),
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
	err := registry.GetRegistryStatus(uri)
	if err != nil {
		utils.Log.Error("Could not add registry " + uri)
	}
	c.Ctx.Redirect(302, "/registries")
}

// ListRegistries returns all registries
func (c *RegistriesController) ListRegistries() {
}

// RemoveRegistry removes a registry from the interface
func (c *RegistriesController) RemoveRegistry() {
}
