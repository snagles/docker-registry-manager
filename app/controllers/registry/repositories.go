package registry

import (
	"github.com/astaxie/beego"
	"github.com/sirupsen/logrus"
	"github.com/snagles/docker-registry-manager/app/models"
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
			logrus.Errorf("Could not connect to registry (%s) to get the repository count: %s ", reg.Registry.URL, err.Error())
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
	c.TplName = "allrepositories.tpl"
}
