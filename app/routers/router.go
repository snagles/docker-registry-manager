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
	beego.Router("/registries/all/count", &controllers.RegistriesController{}, "get:GetRegistryCount")

	// Routers for repositories
	beego.Router("/registries/:registryName/repositories/", &controllers.RepositoriesController{}, "get:GetRepositories")
	beego.Router("/registries/all/repositories/count", &controllers.RepositoriesController{}, "get:GetAllRepositoryCount")

	// Routers for tags
	beego.Router("/registries/:registryName/repositories/*/tags", &controllers.TagsController{}, "get:GetTags")
	beego.Router("/registries/:registryName/repositories/*/tags/:tagName/delete", &controllers.TagsController{}, "post:DeleteTags")

	// Routers for images
	beego.Router("/registries/:registryName/repositories/*/tags/:tagName/images", &controllers.ImagesController{}, "get:GetImages")

	// Routers for logs
	beego.Router("/logs", &controllers.LogsController{}, "get:Get")
}
