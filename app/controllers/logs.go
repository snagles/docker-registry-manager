package controllers

import (
	"bufio"
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/astaxie/beego"
	"github.com/snagles/docker-registry-manager/app/conf"
)

type logResponse struct {
	Error   error  `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

type LogsController struct {
	beego.Controller
}

func (c *LogsController) Get() {
	c.Data["json"] = parseLogs()
	c.ServeJSON()
}

func (c *LogsController) Delete() {
	if err := deleteLog(); err != nil {
		c.Data["json"] = logResponse{
			Error:   err,
			Message: "Failed to clear log",
		}
	} else {
		logrus.Warn("Cleared log file.")
		c.Data["json"] = logResponse{
			Message: "Cleared log.",
		}
	}
	c.ServeJSON()
}

func (c *LogsController) Archive() {
	if err := archiveLog(); err != nil {
		c.Data["json"] = logResponse{
			Error:   err,
			Message: "Failed to archive log",
		}
	} else {
		logrus.Warn("Archived log file.")
		c.Data["json"] = logResponse{
			Message: "Archived log",
		}
	}
	c.ServeJSON()
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
		c.Data["json"] = logResponse{
			Message: "Failed to archive log",
		}
		c.ServeJSON()
		return
	}
	logrus.Warn("Changed log level to " + level)
	c.Data["json"] = logResponse{
		Message: "Changed log level to " + level,
	}
	c.ServeJSON()
}

func (c *LogsController) GetLevel() {
	c.Data["json"] = map[string]interface{}{"log_level": logrus.GetLevel().String(), "log_level_int": logrus.GetLevel()}
	c.ServeJSON()
}

// parseLogs parses the locally stored flat log file that was logged to by logrus
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
