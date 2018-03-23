package registry

import (
	"fmt"
	"net/url"

	"github.com/astaxie/beego"
	"github.com/sirupsen/logrus"
	"github.com/snagles/docker-registry-manager/app/models"
)

// TagsController controls interaction between the UI and all tag related information
type TagsController struct {
	beego.Controller
}

// GetTags returns the template for the registries page
func (c *TagsController) GetTags() {

	registryName := c.Ctx.Input.Param(":registryName")
	repositoryName, _ := url.QueryUnescape(c.Ctx.Input.Param(":splat"))
	repositoryNameEncode := url.QueryEscape(repositoryName)

	manager.AllRegistries.RLock()
	registry, _ := manager.AllRegistries.Registries[registryName]
	repository, _ := registry.Repositories[repositoryName]
	tags := repository.Tags

	c.Data["tags"] = tags
	c.Data["registryName"] = registryName
	c.Data["repositoryNameEncode"] = repositoryNameEncode
	c.Data["repositoryName"] = repositoryName
	manager.AllRegistries.RUnlock()

	// Index template
	c.TplName = "tags.tpl"
}

// DeleteTags deletes the manifest using the passed tag using the digest method
func (c *TagsController) DeleteTags() {
	registryName := c.Ctx.Input.Param(":registryName")
	repositoryName, _ := url.QueryUnescape(c.Ctx.Input.Param(":splat"))
	tag := c.Ctx.Input.Param(":tagName")

	manager.AllRegistries.RLock()
	registry, _ := manager.AllRegistries.Registries[registryName]
	digest, _ := manager.AllRegistries.Registries[registryName].ManifestDigest(repositoryName, tag)
	err := registry.Registry.DeleteManifest(repositoryName, digest)
	manager.AllRegistries.RUnlock()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Digest": digest.String(),
			"Error":  err.Error(),
		}).Errorf("Failed to delete digest.")
		c.CustomAbort(404, fmt.Sprintf("Failure to delete Digest: %v Error: %v", digest.String(), err))
	}
	manager.AllRegistries.Lock()
	ur := registry.Update()
	manager.AllRegistries.Registries[registry.Name] = &ur
	manager.AllRegistries.Unlock()

	c.CustomAbort(200, "Success")

}
