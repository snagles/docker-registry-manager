package controllers

import (
	"github.com/astaxie/beego"
	"github.com/stefannaglee/docker-registry-manager/app/models/registry"
)

type TagsController struct {
	beego.Controller
}

// GetTags returns the template for the registries page
func (c *TagsController) GetTags() {

	registryName := c.Ctx.Input.Param(":registryName")
	repositoryName := c.Ctx.Input.Param(":splat")

	tags, _ := registry.GetTags(registryName, repositoryName)
	c.Data["tags"] = tags.Tags
	c.Data["registryName"] = registryName
	c.Data["repositoryName"] = repositoryName

	// Index template
	c.TplName = "tags.tpl"
}
