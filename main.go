package main

import (
	"flag"

	_ "github.com/DemonVex/docker-registry-manager/routers"
	"github.com/DemonVex/docker-registry-manager/utilities"
	"github.com/Sirupsen/logrus"
	"github.com/astaxie/beego"
)

var logLevel int

func init() {

	// Set and parse the command line flags
	flag.IntVar(&logLevel, "verbosity", 5, "Execution log level of the program: 1 = Panic Level, 2 = Fatal Level, 3 = Error Level, 4 = Warn Level, 5 = Info Level, 6 = Debug Level")
	flag.Parse()

	// Set the log level of the program
	switch {
	case logLevel == 1:
		utils.Log.Level = logrus.PanicLevel
	case logLevel == 2:
		utils.Log.Level = logrus.FatalLevel
	case logLevel == 3:
		utils.Log.Level = logrus.ErrorLevel
	case logLevel == 4:
		utils.Log.Level = logrus.WarnLevel
	case logLevel == 5:
		utils.Log.Level = logrus.InfoLevel
	case logLevel == 6:
		utils.Log.Level = logrus.DebugLevel
	}
}

func main() {

	beego.Run()

}
