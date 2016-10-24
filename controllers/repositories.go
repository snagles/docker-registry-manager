package controllers

import (
	"github.com/astaxie/beego"
	"github.com/snagles/docker-registry-manager/models/registry"
)

type RepositoriesController struct {
	beego.Controller
}

// Get returns the template for the registries page
func (c *RepositoriesController) GetRepositories() {

	registryName := c.Ctx.Input.Param(":registryName")

	if r, ok := registry.Registries[registryName]; ok {

		repositories := r.Repositories

		c.Data["registryName"] = registryName
		c.Data["repositories"] = repositories
		// Index template
		c.TplName = "repositories.tpl"

	}

}
