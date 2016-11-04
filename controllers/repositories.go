package controllers

import (
	"github.com/astaxie/beego"
	"github.com/snagles/docker-registry-manager/models/client"
	"github.com/snagles/docker-registry-manager/models/manager"
	"github.com/snagles/docker-registry-manager/utilities"
)

type RepositoriesController struct {
	beego.Controller
}

// Get returns the template for the registries page
func (c *RepositoriesController) GetRepositories() {

	registryName := c.Ctx.Input.Param(":registryName")

	if r, ok := registry.Registries[registryName]; ok {
		c.Data["registryName"] = registryName
		c.Data["repositories"] = r.Repositories
		// Index template
		c.TplName = "repositories.tpl"
	}
}

func (c *RepositoriesController) GetAllRepositoryCount() {
	c.Data["registries"] = registry.Registries

	var count int
	for _, reg := range registry.Registries {
		repositories, err := client.GetRepositories(reg.URI())
		if err != nil {
			utils.Log.Error("Could not connect to registry to get the repository count. " + err.Error())
		}
		count += len(repositories)
	}
	repositoryCount := struct {
		Count int
	}{
		count,
	}

	c.Data["json"] = &repositoryCount
	c.ServeJSON()
}
