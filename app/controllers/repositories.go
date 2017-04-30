package controllers

import (
	"github.com/Sirupsen/logrus"
	"github.com/astaxie/beego"
	"github.com/snagles/docker-registry-manager/app/models/manager"
)

type RepositoriesController struct {
	beego.Controller
}

// Get returns the template for the registries page
func (c *RepositoriesController) GetRepositories() {

	registryName := c.Ctx.Input.Param(":registryName")

	if r, ok := manager.AllRegistries.Registries[registryName]; ok {
		c.Data["registryName"] = registryName
		c.Data["repositories"] = r.Repositories
	}
	// Index template
	c.TplName = "repositories.tpl"
}

func (c *RepositoriesController) GetAllRepositoryCount() {
	c.Data["registries"] = manager.AllRegistries.Registries

	var count int
	for _, reg := range manager.AllRegistries.Registries {
		repositories, err := reg.Registry.Repositories()
		if err != nil {
			logrus.Error("Could not connect to registry to get the repository count. " + err.Error())
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

	repos := make(map[string][]*manager.Repository)

	for registryName, registry := range manager.AllRegistries.Registries {
		for _, repository := range registry.Repositories {
			repos[registryName] = append(repos[registryName], repository)
		}
	}

	c.Data["repositories"] = repos

	// Index template
	c.TplName = "all_repositories.tpl"
}
