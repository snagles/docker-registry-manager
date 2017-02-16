package controllers

import (
	"github.com/DemonVex/docker-registry-manager/models/client"
	"github.com/DemonVex/docker-registry-manager/models/manager"
	"github.com/DemonVex/docker-registry-manager/utilities"
	"github.com/astaxie/beego"
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

// GetAllRepositories returns the template for the all registries page
func (c *RepositoriesController) GetAllRepositories() {

	repos := make(map[string][]*registry.Repository)

	for registryName, registry := range registry.Registries {
		for _, repository := range registry.Repositories {
			repos[registryName] = append(repos[registryName], repository)
		}
	}

	c.Data["repositories"] = repos

	// Index template
	c.TplName = "all_repositories.tpl"
}
