package controllers

import (
	"encoding/json"
	"fmt"

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

	c.Data["logs"] = logs
	log, _ := json.MarshalIndent(logs, "", "    ")
	fmt.Println(string(log))
	c.ServeJSON()
}
