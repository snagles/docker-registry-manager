package controllers

import (
	"github.com/astaxie/beego"
	"github.com/stefannaglee/docker-registry-manager/app/models/registry"
)

type RepositoriesController struct {
	beego.Controller
}

// Get returns the template for the registries page
func (c *RepositoriesController) GetRepositories() {

	registryName := c.Ctx.Input.Param(":registryName")

	repositories := registry.GetRepositories(registryName)
	c.Data["registryName"] = registryName
	c.Data["repositories"] = repositories
	// Index template
	c.TplName = "repositories.tpl"
}
