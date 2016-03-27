package routers

import (
	"github.com/astaxie/beego"
	"github.com/stefannaglee/docker-registry-manager/app/controllers"
)

func init() {
	beego.Router("/", &controllers.RegistriesController{})

	// Routers for registries
	beego.Router("/registries", &controllers.RegistriesController{})
	beego.Router("/registries/", &controllers.RegistriesController{})
	beego.Router("/registries/:registryName/repositories/", &controllers.RepositoriesController{}, "get:GetRepositories")
	beego.Router("/registries/:registryName/repositories/:repositoryName/tags", &controllers.TagsController{}, "get:GetTags")
}
