package registry

import (
	"github.com/astaxie/beego"
	"github.com/snagles/docker-registry-manager/app/models"
)

// RepositoriesController controls interactions between the UI and any repository related information
type RepositoriesController struct {
	beego.Controller
}

// GetRepositories returns the template for the table on the repositories page
func (c *RepositoriesController) GetRepositories() {

	registryName := c.Ctx.Input.Param(":registryName")

	if r, ok := manager.AllRegistries.Registries[registryName]; ok {
		c.Data["registryName"] = registryName
		c.Data["repositories"] = r.Repositories
	}
	// Index template
	c.TplName = "repositories.tpl"
}

// GetAllRepositoryCount returns the number of currently available repositories
func (c *RepositoriesController) GetAllRepositoryCount() {
	c.Data["registries"] = manager.AllRegistries.Registries

	var repositoryCount struct {
		Count int
	}
	for _, reg := range manager.AllRegistries.Registries {
		repositoryCount.Count += len(reg.Repositories)
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
