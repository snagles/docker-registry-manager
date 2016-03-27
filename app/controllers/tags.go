package controllers

import (
	"github.com/astaxie/beego"
	"github.com/stefannaglee/docker-registry-manager/app/models/registry"
)

type TagsController struct {
	beego.Controller
}

// Get returns the template for the registries page
func (c *TagsController) GetTags() {

	registryName := c.Ctx.Input.Param(":registryName")
	repositoryName := c.Ctx.Input.Param(":repositoryName")

	tags, _ := registry.GetTags(registryName, repositoryName)
	c.Data["tags"] = tags.Tags

	// Index template
	c.TplName = "tags.tpl"
}
