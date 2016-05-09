package controllers

import (
	"github.com/astaxie/beego"
	"github.com/stefannaglee/docker-registry-manager/app/utilities"
)

type SettingsController struct {
	beego.Controller
}

func (c *SettingsController) Get() {

	logs := utils.ParseLogFile()

	c.Data["logs"] = logs

	c.TplName = "settings.tpl"
}
