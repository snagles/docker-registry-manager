package controllers

import (
	"net/url"

	"github.com/astaxie/beego"
	"github.com/stefannaglee/docker-registry-manager/app/models/registry"
)

type TagsController struct {
	beego.Controller
}

// GetTags returns the template for the registries page
func (c *TagsController) GetTags() {

	registryName := c.Ctx.Input.Param(":registryName")
	repositoryName, _ := url.QueryUnescape(c.Ctx.Input.Param(":splat"))
	repositoryNameEncode := url.QueryEscape(repositoryName)

	tags, _ := registry.GetTagsForView(registryName, repositoryName)

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

	registry.DeleteTag(registryName, repositoryName, tag)

	c.CustomAbort(200, "Success")

}
