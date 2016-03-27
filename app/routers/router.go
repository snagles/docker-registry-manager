package routers

import (
	"github.com/astaxie/beego"
	"github.com/stefannaglee/docker-registry-manager/app/controllers"
)

func init() {
	beego.Router("/", &controllers.RegistriesController{})
}
