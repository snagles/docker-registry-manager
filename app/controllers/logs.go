package controllers

import (
	"github.com/astaxie/beego"
	"github.com/stefannaglee/docker-registry-manager/app/utilities"
)

type LogsController struct {
	beego.Controller
}

func (c *LogsController) Get() {

	logs := utils.ParseLogFile()

	c.Data["logs"] = logs

	c.TplName = "logs.tpl"
}
