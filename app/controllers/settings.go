package controllers

import (
	"github.com/astaxie/beego"
	"github.com/stefannaglee/docker-registry-manager/app/utilities"
)

type SettingsController struct {
	beego.Controller
}

func (c *SettingsController) Get() {

	c.TplName = "settings.tpl"
}

func (c *SettingsController) GetLogs() {
	logs := utils.ParseLogFile()

	c.Data["json"] = &logs
	c.ServeJSON()
}
