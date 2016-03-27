package routers

import (
	"github.com/stefannaglee/docker-registry-manager/app/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
}
