package controllers

import (
	"github.com/Sirupsen/logrus"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/toolbox"
	"github.com/snagles/docker-registry-manager/settings"
)

type SettingsController struct {
	beego.Controller
}

func (c *SettingsController) Get() {
	c.Data["releaseVersion"] = settings.ReleaseVersion
	c.Data["activeLevel"] = logrus.GetLevel()
	c.Data["allLevels"] = logrus.AllLevels
	c.TplName = "settings.tpl"
}

func (c *SettingsController) GetReleaseVersion() {
	currentRelease := settings.ReleaseVersion
	type ReleaseVersion struct {
		ReleaseVersion string
	}
	r := ReleaseVersion{
		ReleaseVersion: currentRelease,
	}
	c.Data["json"] = &r
	c.ServeJSON()
}

// GetLiveStatistics returns stats on request information tracked by beego
func (c *SettingsController) GetLiveStatistics() {

	r := toolbox.StatisticsMap
	rs := r.GetMapData()

	// Convert beego times to seconds for sorted
	for _, req := range rs {
		var err error
		req["min_s"], err = settings.StatToSeconds(req["min_time"].(string))
		if err != nil {
			logrus.Error(err)
		}
		req["max_s"], err = settings.StatToSeconds(req["max_time"].(string))
		if err != nil {
			logrus.Error(err)
		}
		req["avg_s"], err = settings.StatToSeconds(req["avg_time"].(string))
		if err != nil {
			logrus.Error(err)
		}
		req["total_s"], err = settings.StatToSeconds(req["total_time"].(string))
		if err != nil {
			logrus.Error(err)
		}
	}

	c.Data["json"] = &rs
	c.ServeJSON()
}