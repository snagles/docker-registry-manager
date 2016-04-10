package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/astaxie/beego"
	"github.com/stefannaglee/docker-registry-manager/app/models/registry"
	_ "github.com/stefannaglee/docker-registry-manager/app/routers"
	"github.com/stefannaglee/docker-registry-manager/app/utilities"
)

// DuplicateFlags slice contains a list of registries to use.
type DuplicateFlags []string

// String is the method to format the flag's value, part of the flag.Value interface.
func (d *DuplicateFlags) String() string {
	return fmt.Sprint(*d)
}

// Set builds the registry slice by taking in the comma separated list and splitting on commas
func (d *DuplicateFlags) Set(value string) error {
	for _, url := range strings.Split(value, ",") {
		*d = append(*d, url)
	}
	return nil
}

var logLevel int
var registryFlags DuplicateFlags

func init() {

	// Set and parse the command line flags
	flag.IntVar(&logLevel, "verbosity", 5, "Execution log level of the program: 1 = Panic Level, 2 = Fatal Level, 3 = Error Level, 4 = Warn Level, 5 = Info Level, 6 = Debug Level")
	flag.Var(&registryFlags, "registry", "comma-separated list of registries to use. e.g https://host.domain:5000/v2/")
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
	utils.Log.Error("Hi")
}

func main() {

	// Loop through the slice of passed registries and test their status
	for _, regString := range registryFlags {
		if err := registry.GetRegistryStatus(regString); err != nil {
			// Notify of success
			utils.Log.WithFields(logrus.Fields{
				"Error": err,
			}).Fatal("We are unable to make a connection to the registry!")

			// Exit the program
			os.Exit(1)
		}
	}

	beego.Run()

}
