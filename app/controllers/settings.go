package controllers

import "github.com/astaxie/beego"

type SettingsController struct {
	beego.Controller
}

func (c *SettingsController) Get() {

	c.TplName = "settings.tpl"
}
