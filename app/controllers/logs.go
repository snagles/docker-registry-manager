package controllers

import (
	"bufio"
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/astaxie/beego"
	"github.com/snagles/docker-registry-manager/app/conf"
)

type LogsController struct {
	beego.Controller
}

func (c *LogsController) Get() {
	c.Data["json"] = parseLogs()
	c.ServeJSON()
}

func (c *LogsController) Delete() {
	if err := deleteLog(); err != nil {
		c.CustomAbort(404, "Failed to clear log: "+err.Error())
	}
	logrus.Warn("Cleared log file.")
	c.CustomAbort(200, "Success")
}

func (c *LogsController) Archive() {
	if err := archiveLog(); err != nil {
		c.CustomAbort(404, "Failed to clear log: "+err.Error())
	}
	logrus.Warn("Archived log file.")
	c.CustomAbort(200, "Success")
}

func (c *LogsController) PostLevel() {
	level := c.Ctx.Input.Param(":level")
	switch {
	case level == "panic":
		logrus.SetLevel(logrus.PanicLevel)
	case level == "fatal":
		logrus.SetLevel(logrus.FatalLevel)
	case level == "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case level == "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case level == "info":
		logrus.SetLevel(logrus.InfoLevel)
	case level == "debug":
		logrus.SetLevel(logrus.DebugLevel)
	default:
		c.CustomAbort(404, "Unrecognized log level")
	}
	c.CustomAbort(200, "Success")
}

func (c *LogsController) GetLevel() {
	c.Data["json"] = map[string]interface{}{"log_level": logrus.GetLevel().String(), "log_level_int": logrus.GetLevel()}
	c.ServeJSON()
}

// parseLogs parseLogss the locally stored flat log file that was logged to by logrus
//{"file":"log.go","level":"warning","line":588,"msg":"test","source":"beego","time":"2017-04-29T20:37:09-04:00"}

type Entry struct {
	File    string    `json:"file"`
	Level   string    `json:"level"`
	Line    int       `json:"line"`
	Message string    `json:"msg"`
	Source  string    `json:"source"`
	Time    time.Time `json:"time"`
}

func parseLogs() []Entry {
	file, err := os.Open(conf.LogFile)
	defer file.Close()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err.Error(),
		}).Error("Failed to open configured log file: " + conf.LogFile)
	}

	es := []Entry{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		e := Entry{}
		err := json.Unmarshal(scanner.Bytes(), &e)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"Error": err.Error(),
			}).Error("Failed to parse log file line")
			return es
		}
		es = append(es, e)
	}

	if err := scanner.Err(); err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err.Error(),
		}).Error("Failed to parseLogs log file line")
	}
	return es
}

// deleteLogs empties the log file
func deleteLog() error {
	// Create the new file
	if _, err := os.Create(conf.LogFile + ".new"); err != nil {
		logrus.Error(err)
		return err
	}

	// Rename the old error log, and update the name of the new one
	if err := os.Rename(conf.LogFile, conf.LogFile+".old"); err != nil {
		logrus.Error(err)
		return err
	}
	if err := os.Rename(conf.LogFile+".new", conf.LogFile); err != nil {
		logrus.Error(err)
		return err
	}

	if err := os.Remove(conf.LogFile + ".old"); err != nil {
		logrus.Error(err)
	}
	f, _ := os.OpenFile(conf.LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	logrus.SetOutput(f)
	return nil
}

// archiveLog creates a backup of the current log in the logs directory, creates a new log file, and writes to the new log
func archiveLog() error {
	// Create the new file
	if _, err := os.Create(conf.LogFile + ".new"); err != nil {
		logrus.Error(err)
		return err
	}

	// Rename the old error log, and update the name of the new one
	if err := os.Rename(conf.LogFile, conf.LogDir+strconv.Itoa(int(time.Now().Unix()))+".json"); err != nil {
		logrus.Error(err)
		return err
	}
	if err := os.Rename(conf.LogFile+".new", conf.LogFile); err != nil {
		logrus.Error(err)
		return err
	}

	f, _ := os.OpenFile(conf.LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	logrus.SetOutput(f)
	return nil
}
