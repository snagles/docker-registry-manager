package controllers

import (
	"github.com/astaxie/beego"
	"github.com/stefannaglee/docker-registry-manager/app/models/registry"
)

type ImagesController struct {
	beego.Controller
}

// GetImages returns the template for the images page
func (c *ImagesController) GetImages() {

	registryName := c.Ctx.Input.Param(":registryName")
	repositoryName := c.Ctx.Input.Param(":splat")
	tagName := c.Ctx.Input.Param(":tagName")

	img, _ := registry.GetImage(registryName, repositoryName, tagName)

	c.Data["history"] = img.History
	c.Data["registryName"] = registryName
	c.Data["repositoryName"] = repositoryName
	c.Data["tagName"] = tagName

	// Index template
	c.TplName = "images.tpl"
}
