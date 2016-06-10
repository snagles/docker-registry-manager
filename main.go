package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/astaxie/beego"
	"github.com/stefannaglee/docker-registry-manager/models/registry"
	_ "github.com/stefannaglee/docker-registry-manager/routers"
	"github.com/stefannaglee/docker-registry-manager/utilities"
)

const (
	// Update version for each release
	MajorVersion = 1
	MinorVersion = 0
	PatchVersion = 0
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

	utils.ReleaseVersion = strconv.Itoa(MajorVersion) + "." + strconv.Itoa(MinorVersion) + "." + strconv.Itoa(PatchVersion)

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
}

func main() {

	// Loop through the slice of passed registries and test their status
	for _, regString := range registryFlags {

		// Build the registry object
		r, err := registry.ParseRegistry(regString)
		if err != nil {
			// Notify of failure to parse
			utils.Log.WithFields(logrus.Fields{
				"Error": err,
			}).Fatal("We are unable to determine the URI for the registry!")

			// Exit the program
			os.Exit(1)
		}

		// Check to see if the registry is available
		err = r.UpdateRegistryStatus()
		if err != nil && r.Status != "available" {
			// Notify of success
			utils.Log.WithFields(logrus.Fields{
				"Error": err,
			}).Fatal("We are unable to connection to the registry!")

			// Exit the program
			os.Exit(1)
		}

		// Add the registry to the map of active registries
		r.AddRegistry()
	}

	beego.Run()

}
