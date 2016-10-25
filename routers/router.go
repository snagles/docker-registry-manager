package routers

import (
	"github.com/astaxie/beego"
	"github.com/snagles/docker-registry-manager/controllers"
)

func init() {
	beego.Router("/", &controllers.RegistriesController{})
	beego.Router("/.release", &controllers.SettingsController{})

	// Routers for registries
	beego.Router("/registries", &controllers.RegistriesController{})
	beego.Router("/registries/", &controllers.RegistriesController{})
	beego.Router("/registries/all/count", &controllers.RegistriesController{}, "get:GetRegistryCount")
	beego.Router("/registries/add", &controllers.RegistriesController{}, "post:AddRegistry")

	// Routers for repositories
	beego.Router("/registries/:registryName/repositories/", &controllers.RepositoriesController{}, "get:GetRepositories")

	// Routers for tags
	beego.Router("/registries/:registryName/repositories/*/tags", &controllers.TagsController{}, "get:GetTags")

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
