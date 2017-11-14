package conf

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	beego "github.com/astaxie/beego/logs"
	. "github.com/smartystreets/goconvey/convey"
)

// TestGoPath tests that the GOPATH is set correctly
func TestGoPath(t *testing.T) {
	Convey("We should be able to retrieve the GOPATH", t, func(c C) {
		c.So(GOPATH, ShouldNotBeBlank)
	})
}

// TestLogFileCreation tests that the logfile is successfully created
func TestLogFileCreation(t *testing.T) {
	Convey("The log file and directory should be set correctly", t, func(c C) {
		c.So(LogDir, ShouldNotBeBlank)
		c.So(LogFile, ShouldNotBeBlank)
	})

	_, err := os.Stat(LogFile)
	Convey("The log file and directory should be created successfully", t, func(c C) {
		c.So(err, ShouldBeNil)
	})
}

// TestBeegoLog tests that the GOPATH is set correctly
func TestBeegoLog(t *testing.T) {
	beego.Error("%s", "test")
	logrus.Error("test")
	Convey("We should be able to retrieve the GOPATH", t, func(c C) {
		c.So(GOPATH, ShouldNotBeBlank)
	})
}
