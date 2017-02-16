package routers

import (
	"github.com/DemonVex/docker-registry-manager/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.RegistriesController{})
	beego.Router("/.release", &controllers.SettingsController{})

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
	beego.Router("/logs", &controllers.SettingsController{}, "get:GetLogs")
	beego.Router("/logs/clear", &controllers.SettingsController{}, "post:ClearLogs")
	beego.Router("/logs/archive", &controllers.SettingsController{}, "post:ArchiveLogs")
	beego.Router("/logs/toggle-debug", &controllers.SettingsController{}, "get:ToggleDebug")
	beego.Router("/logs/level", &controllers.SettingsController{}, "get:GetLogLevel")

	// Routers for settings
	beego.Router("/settings", &controllers.SettingsController{}, "get:Get")
	beego.Router("/settings/stats", &controllers.SettingsController{}, "get:GetLiveStatistics")
}
