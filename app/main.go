package main

import (
	"flag"

	"github.com/Sirupsen/logrus"
	"github.com/astaxie/beego"
	_ "github.com/snagles/docker-registry-manager/app/routers"
)

var logLevel int

func init() {

	// Set and parse the command line flags
	flag.IntVar(&logLevel, "verbosity", 5, "Execution log level of the program: 1 = Panic Level, 2 = Fatal Level, 3 = Error Level, 4 = Warn Level, 5 = Info Level, 6 = Debug Level")
	flag.Parse()

	// Set the log level of the program
	switch {
	case logLevel == 1:
		logrus.SetLevel(logrus.PanicLevel)
	case logLevel == 2:
		logrus.SetLevel(logrus.FatalLevel)
	case logLevel == 3:
		logrus.SetLevel(logrus.ErrorLevel)
	case logLevel == 4:
		logrus.SetLevel(logrus.WarnLevel)
	case logLevel == 5:
		logrus.SetLevel(logrus.InfoLevel)
	case logLevel == 6:
		logrus.SetLevel(logrus.DebugLevel)
	}
}

func main() {
	beego.Run()
}
