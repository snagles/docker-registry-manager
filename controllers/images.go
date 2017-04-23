package controllers

import (
	"net/url"

	"github.com/astaxie/beego"
	"github.com/snagles/docker-registry-manager/models/manager"
)

type ImagesController struct {
	beego.Controller
}

// GetImages returns the template for the images page
func (c *ImagesController) GetImages() {

	registryName := c.Ctx.Input.Param(":registryName")
	repositoryName, _ := url.QueryUnescape(c.Ctx.Input.Param(":splat"))
	repositoryNameEncode := url.QueryEscape(repositoryName)
	c.Data["tagName"] = c.Ctx.Input.Param(":tagName")

	registry, _ := manager.AllRegistries.Registries[registryName]
	c.Data["registry"] = registry

	tag, _ := registry.Repositories[repositoryName].Tags[c.Ctx.Input.Param(":tagName")]
	c.Data["tag"] = tag

	c.Data["registryName"] = registryName
	c.Data["repositoryNameEncode"] = repositoryNameEncode
	c.Data["repositoryName"] = repositoryName

	// Index template
	c.TplName = "images.tpl"
}
