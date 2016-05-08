package routers

import (
	"github.com/astaxie/beego"
	"github.com/stefannaglee/docker-registry-manager/controllers"
)

func init() {
	beego.Router("/", &controllers.RegistriesController{})

	// Routers for registries
	beego.Router("/registries", &controllers.RegistriesController{})
	beego.Router("/registries/", &controllers.RegistriesController{})
	beego.Router("/registries/all/count", &controllers.RegistriesController{}, "get:GetRegistryCount")
	beego.Router("/registries/add", &controllers.RegistriesController{}, "post:AddRegistry")
	beego.Router("/registries/test", &controllers.RegistriesController{}, "post:TestRegistryStatus")

	// Routers for repositories
	beego.Router("/registries/:registryName/repositories/", &controllers.RepositoriesController{}, "get:GetRepositories")
	beego.Router("/registries/all/repositories/count", &controllers.RepositoriesController{}, "get:GetAllRepositoryCount")
	beego.Router("/registries/all/repositories", &controllers.RepositoriesController{}, "get:GetAllRepositories")

	// Routers for tags
	beego.Router("/registries/:registryName/repositories/*/tags", &controllers.TagsController{}, "get:GetTags")
	beego.Router("/registries/:registryName/repositories/*/tags/:tagName/delete", &controllers.TagsController{}, "post:DeleteTags")

	// Routers for images
	beego.Router("/registries/:registryName/repositories/*/tags/:tagName/images", &controllers.ImagesController{}, "get:GetImages")

	// Routers for logs
	beego.Router("/logs", &controllers.LogsController{}, "get:Get")

	// Routers for settings
	beego.Router("/settings", &controllers.SettingsController{}, "get:Get")
}
