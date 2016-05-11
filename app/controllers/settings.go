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

func (c *SettingsController) ClearLogs() {

	err := utils.ClearLogFile()
	if err == nil {
		c.CustomAbort(200, "Success")
	}
	c.CustomAbort(404, "Failed to clear log: "+err.Error())

}
