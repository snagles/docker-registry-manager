package routers

import (
	"github.com/astaxie/beego"
	"github.com/snagles/docker-registry-manager/app/controllers"
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
	beego.Router("/registries/:registryName/repositories", &controllers.RepositoriesController{}, "get:GetRepositories")
	beego.Router("/registries/all/repositories/count", &controllers.RepositoriesController{}, "get:GetAllRepositoryCount")
	beego.Router("/registries/all/repositories", &controllers.RepositoriesController{}, "get:GetAllRepositories")

	// Routers for tags
	beego.Router("/registries/:registryName/repositories/*/tags", &controllers.TagsController{}, "get:GetTags")
	beego.Router("/registries/:registryName/repositories/*/tags/:tagName/delete", &controllers.TagsController{}, "post:DeleteTags")

	// Routers for images
	beego.Router("/registries/:registryName/repositories/*/tags/:tagName/images", &controllers.ImagesController{}, "get:GetImages")

	// Routers for logs
	beego.Router("/logs", &controllers.LogsController{}, "get:Get")
	beego.Router("/logs", &controllers.LogsController{}, "delete:Delete")
	beego.Router("/logs/archive", &controllers.LogsController{}, "post:Archive")
	beego.Router("/logs/level", &controllers.LogsController{}, "get:GetLevel")
	beego.Router("/logs/level/:level", &controllers.LogsController{}, "post:PostLevel")

	// Routers for events
	beego.Router("/envelope", &controllers.EventsController{}, "post:PostEnvelope")
	beego.Router("/events", &controllers.EventsController{}, "get:Get")
	beego.Router("/events/:registryName", &controllers.EventsController{}, "get:GetRegistryEvents")
	beego.Router("/events/:registryName/:eventID", &controllers.EventsController{}, "get:GetRegistryEventID")

	// Routers for settings
	beego.Router("/settings", &controllers.SettingsController{}, "get:Get")
	beego.Router("/settings/stats", &controllers.SettingsController{}, "get:GetLiveStatistics")
}
