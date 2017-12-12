package routers

import (
	"github.com/astaxie/beego"
	"github.com/snagles/docker-registry-manager/app/controllers/app"
	"github.com/snagles/docker-registry-manager/app/controllers/registry"
)

func init() {
	beego.Router("/", &registry.RegistriesController{})

	/* proposed
	beego.Router("/logs/requests", &app.SettingsController{}, "get:GetLiveStatistics")
	beego.Router("/logs/active-level", &app.LogsController{}, "get:GetLevel")
	beego.Router("/logs/actions/archive", &app.LogsController{}, "post:Archive")
	beego.Router("/logs/actions/delete", &app.LogsController{}, "delete:Delete")
	beego.Router("/logs/actions/set-level/:level", &app.LogsController{}, "post:PostLevel")
	*/

	// Routers for settings
	beego.Router("/settings/stats", &app.SettingsController{}, "get:GetLiveStatistics")

	// Routers for logs
	beego.Router("/logs", &app.LogsController{}, "get:Get")
	beego.Router("/logs", &app.LogsController{}, "delete:Delete")
	beego.Router("/logs/archive", &app.LogsController{}, "post:Archive")
	beego.Router("/logs/level", &app.LogsController{}, "get:GetLevel")
	beego.Router("/logs/level/:level", &app.LogsController{}, "post:PostLevel")

	// Routers for events
	beego.Router("/envelope", &app.EventsController{}, "post:PostEnvelope")
	beego.Router("/events", &app.EventsController{}, "get:Get")
	beego.Router("/events/:registryName", &app.EventsController{}, "get:GetRegistryEvents")
	beego.Router("/events/:registryName/:eventID", &app.EventsController{}, "get:GetRegistryEventID")

	// Routers for registries
	beego.Router("/registries", &registry.RegistriesController{})
	beego.Router("/registries/", &registry.RegistriesController{})
	beego.Router("/registries/all/count", &registry.RegistriesController{}, "get:GetRegistryCount")
	beego.Router("/registries/add", &registry.RegistriesController{}, "post:AddRegistry")
	beego.Router("/registries/test", &registry.RegistriesController{}, "post:TestRegistryStatus")

	// Routers for repositories
	beego.Router("/registries/:registryName/repositories", &registry.RepositoriesController{}, "get:GetRepositories")
	beego.Router("/registries/all/repositories/count", &registry.RepositoriesController{}, "get:GetAllRepositoryCount")
	beego.Router("/registries/all/repositories", &registry.RepositoriesController{}, "get:GetAllRepositories")

	// Routers for tags
	beego.Router("/registries/:registryName/repositories/*/tags", &registry.TagsController{}, "get:GetTags")
	beego.Router("/registries/:registryName/repositories/*/tags/:tagName/delete", &registry.TagsController{}, "post:DeleteTags")

	// Routers for images
	beego.Router("/registries/:registryName/repositories/*/tags/:tagName/images", &registry.ImagesController{}, "get:GetImages")

}
