package controllers

import (
	"net/url"

	"github.com/DemonVex/docker-registry-manager/models/client"
	"github.com/DemonVex/docker-registry-manager/models/manager"
	"github.com/astaxie/beego"
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

func (c *TagsController) DeleteTags() {
	registryName := c.Ctx.Input.Param(":registryName")
	repositoryName, _ := url.QueryUnescape(c.Ctx.Input.Param(":splat"))
	tag := c.Ctx.Input.Param(":tagName")

	registry, _ := registry.Registries[registryName]
	success, err := client.DeleteTag(registry.URI(), repositoryName, tag)
	if !success || err != nil {
		c.CustomAbort(404, "Failure")
	}

	c.CustomAbort(200, "Success")

}
