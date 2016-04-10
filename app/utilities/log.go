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
		logrus.InfoLevel:  "./logs/info.log",
		logrus.WarnLevel:  "./logs/warn.log",
		logrus.FatalLevel: "./logs/fatal.log",
	}))
}
