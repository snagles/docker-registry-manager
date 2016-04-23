package controllers

import (
	"github.com/astaxie/beego"
	"github.com/stefannaglee/docker-registry-manager/app/models/registry"
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

// ListRegistries returns all registries
func (c *RegistriesController) ListRegistries() {
}

// RemoveRegistry removes a registry from the interface
func (c *RegistriesController) RemoveRegistry() {
}
