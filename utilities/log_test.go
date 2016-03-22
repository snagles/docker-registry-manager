package utils

import (
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
)

// TestParseLog tests parsing of the local log file
func TestParseLog(t *testing.T) {

	// Create an example log
	Log.Error("test")

	logs := ParseLogFile()
	Convey("We should be able to parse the log file ", t, func() {
		So(len(logs), ShouldBeGreaterThan, 0)
	})

}

// TestClearLogFile tests clearing of the local log file
func TestClearLogFile(t *testing.T) {

	err := ClearLogFile()
	Convey("We should be able to clear the log file ", t, func() {
		So(err, ShouldBeNil)
	})

	// Check the log file for our clear message
	contents, _ := ioutil.ReadFile(LogFile)
	r, _ := regexp.Compile("Truncated and cleared the log file!")
	contains := r.MatchString(string(contents))
	Convey("The log should show that the log was cleared", t, func() {
		So(contains, ShouldBeTrue)
	})

}

// TestArchiveLogFile tests archiving of the local log file
func TestArchiveLogFile(t *testing.T) {

	err := ArchiveLogFile()
	Convey("We should be able archive the log file ", t, func() {
		So(err, ShouldBeNil)
	})

	logTime := time.Now().Add(-2 * time.Second)
	var exists bool
	// Give it a 5 second buffer and see if there is a log file somewhere within that range
	for i := 0; i < 5; i++ {
		logTime = logTime.Add(1 * time.Second)
		_, err = os.Stat(LogPath + strconv.Itoa(int(logTime.Unix())) + "-error.log")
		if err == nil {
			exists = true
			break
		}
	}
	Convey("The archived file should be in the log directory", t, func() {
		So(exists, ShouldBeTrue)
	})

}

// TestToggleDebug tests the toggling of debug mode
func TestToggleDebug(t *testing.T) {

	currentLevel := Log.Level
	ToggleDebug()
	if currentLevel == logrus.DebugLevel {
		Convey("The log level should be toggled from debug to info", t, func() {
			So(Log.Level, ShouldEqual, logrus.InfoLevel)
		})
		ToggleDebug()
		Convey("The log level should be toggled from info to debug", t, func() {
			So(Log.Level, ShouldEqual, logrus.DebugLevel)
		})
	} else {
		Convey("The log level should be toggled from info to debug", t, func() {
			So(Log.Level, ShouldEqual, logrus.DebugLevel)
		})
		ToggleDebug()
		Convey("The log level should be toggled from debug to info", t, func() {
			So(Log.Level, ShouldEqual, logrus.InfoLevel)
		})
	}

}
