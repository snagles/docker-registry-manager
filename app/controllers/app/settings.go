package app

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/toolbox"
	"github.com/sirupsen/logrus"
)

type SettingsController struct {
	beego.Controller
}

func (c *SettingsController) Get() {
	c.Data["activeLevel"] = logrus.GetLevel()
	c.Data["allLevels"] = logrus.AllLevels
	c.Data["logs"] = parseLogs()
	c.Data["stats"] = toolbox.StatisticsMap.GetMapData()
	c.TplName = "settings.tpl"
}

// GetLiveStatistics returns stats on request information tracked by beego
func (c *SettingsController) GetLiveStatistics() {
	c.Data["json"] = toolbox.StatisticsMap.GetMapData()
	c.ServeJSON()
}
