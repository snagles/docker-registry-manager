package controllers

import (
	"net/url"

	"github.com/astaxie/beego"
	"github.com/snagles/docker-registry-manager/models/registry"
)

type TagsController struct {
	beego.Controller
}

// GetTags returns the template for the registries page
func (c *TagsController) GetTags() {

	registryName := c.Ctx.Input.Param(":registryName")
	repositoryName, _ := url.QueryUnescape(c.Ctx.Input.Param(":splat"))
	repositoryNameEncode := url.QueryEscape(repositoryName)

	registry, _ := registry.Registries[registryName]
	repository, _ := registry.Repositories[repositoryName]
	tags := repository.Tags

	c.Data["tags"] = tags
	c.Data["registryName"] = registryName
	c.Data["repositoryNameEncode"] = repositoryNameEncode
	c.Data["repositoryName"] = repositoryName

	// Index template
	c.TplName = "tags.tpl"
}
