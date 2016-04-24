package controllers

import (
	"net/url"

	"github.com/astaxie/beego"
	"github.com/stefannaglee/docker-registry-manager/app/models/registry"
)

type ImagesController struct {
	beego.Controller
}

// GetImages returns the template for the images page
func (c *ImagesController) GetImages() {

	registryName := c.Ctx.Input.Param(":registryName")
	repositoryName, _ := url.QueryUnescape(c.Ctx.Input.Param(":splat"))
	tagName := c.Ctx.Input.Param(":tagName")
	repositoryNameEncode := url.QueryEscape(repositoryName)

	img, _ := registry.GetImage(registryName, repositoryName, tagName)

	tagInfo, _ := registry.GetTag(registryName, repositoryName, tagName)

	c.Data["containsV1Size"] = img.ContainsV1Size
	c.Data["os"] = img.History[0].V1Compatibility.Os
	c.Data["arch"] = img.History[0].V1Compatibility.Architecture
	c.Data["history"] = img.History
	c.Data["registryName"] = registryName
	c.Data["repositoryName"] = repositoryName
	c.Data["repositoryNameEncode"] = repositoryNameEncode
	c.Data["tagInfo"] = tagInfo

	// Index template
	c.TplName = "images.tpl"
}
