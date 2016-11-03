package controllers

import (
	"net/url"

	"github.com/astaxie/beego"
	"github.com/snagles/docker-registry-manager/models/client"
	"github.com/snagles/docker-registry-manager/models/registry"
)

type ImagesController struct {
	beego.Controller
}

// GetImages returns the template for the images page
func (c *ImagesController) GetImages() {

	registryName := c.Ctx.Input.Param(":registryName")
	repositoryName, _ := url.QueryUnescape(c.Ctx.Input.Param(":splat"))

	c.Data["tagName"] = c.Ctx.Input.Param(":tagName")

	repositoryNameEncode := url.QueryEscape(repositoryName)

	registry := registry.Registries[registryName]
	img, _ := client.GetImage(registry.URI(), repositoryName, c.Ctx.Input.Param(":tagName"))

	c.Data["registry"] = registry

	c.Data["image"] = img
	c.Data["containsV1Size"] = img.ContainsV1Size
	c.Data["os"] = img.History[0].V1Compatibility.Os
	c.Data["arch"] = img.History[0].V1Compatibility.Architecture
	c.Data["history"] = img.History
	c.Data["registryName"] = registryName
	c.Data["repositoryName"] = repositoryName
	c.Data["repositoryNameEncode"] = repositoryNameEncode
	c.Data["layers"] = img.FsLayers

	// Index template
	c.TplName = "images.tpl"
}
