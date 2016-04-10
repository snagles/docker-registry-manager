package utils

import (
	"github.com/Sirupsen/logrus"
	"github.com/rifflock/lfshook"
)

var Log = logrus.New()

func init() {

	Log.Formatter = &logrus.JSONFormatter{}
	Log.Hooks.Add(lfshook.NewHook(lfshook.PathMap{
		logrus.ErrorLevel: "./logs/error.log",
		logrus.InfoLevel:  "./app/logs/info.log",
	}))
}
