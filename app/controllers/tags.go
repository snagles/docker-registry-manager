package controllers

import (
	"fmt"
	"net/url"

	"github.com/sirupsen/logrus"
	"github.com/astaxie/beego"
	"github.com/snagles/docker-registry-manager/app/models"
)

type TagsController struct {
	beego.Controller
}

// GetTags returns the template for the registries page
func (c *TagsController) GetTags() {

	registryName := c.Ctx.Input.Param(":registryName")
	repositoryName, _ := url.QueryUnescape(c.Ctx.Input.Param(":splat"))
	repositoryNameEncode := url.QueryEscape(repositoryName)

	registry, _ := manager.AllRegistries.Registries[registryName]
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

	registry, _ := manager.AllRegistries.Registries[registryName]
	digest, _ := manager.AllRegistries.Registries[registryName].ManifestDigest(repositoryName, tag)
	err := registry.Registry.DeleteManifest(repositoryName, digest)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Digest": digest.String(),
			"Error":  err.Error(),
		}).Errorf("Failed to delete digest.")
		c.CustomAbort(404, fmt.Sprintf("Failure to delete Digest: %v Error: %v", digest.String(), err))
	}
	registry.Refresh()

	c.CustomAbort(200, "Success")

}
