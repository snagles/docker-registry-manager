package utils

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/rifflock/lfshook"
)

// GoPath contains the gopath of the env
var GoPath string

// AppPath contains the app path of the env
var AppPath string

// LogPath contains the path to store the logs
var LogPath string

// LogFile contains the fully qualified log file location
var LogFile string

// LogEntry contains information unmarshalled from logrus that was logged to a file
//{"level":"error","msg":"Test","time":"2016-04-10T12:05:30-04:00"}
type LogEntry struct {
	Level   string    `json:"level"`
	Message string    `json:"msg"`
	Time    time.Time `json:"time"`
}

// Log creates a new logrus instance that can be exported and used throughout the project
var Log = logrus.New()

func init() {

	goPath := os.Getenv("GOPATH")
	GoPath = goPath
	appPath := goPath + "/src/github.com/snagles/docker-registry-manager/"
	AppPath = appPath
	logPath := appPath + "/logs/"
	LogPath = logPath
	logFile := logPath + "/error.log"
	LogFile = logFile

	// Create the log directory if needed
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		if err := os.Mkdir(logPath, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}

	// Create the log file if needed
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		if _, err = os.Create(logFile); err != nil {
			log.Fatal(err)
		}
	}

	Log.Formatter = &logrus.JSONFormatter{}
	lfsHook := lfshook.NewHook(lfshook.PathMap{
		logrus.ErrorLevel: logFile,
		logrus.InfoLevel:  logFile,
		logrus.WarnLevel:  logFile,
		logrus.FatalLevel: logFile,
		logrus.PanicLevel: logFile,
		logrus.DebugLevel: logFile,
	})
	lfsHook.SetFormatter(&logrus.JSONFormatter{})
	Log.Hooks.Add(lfsHook)
}

// ParseLogFile parses the locally stored flat log file that was logged to by logrus
func ParseLogFile() []LogEntry {
	file, err := os.Open(LogFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	es := []LogEntry{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		e := LogEntry{}
		err := json.Unmarshal(scanner.Bytes(), &e)
		if err != nil {
			continue
		}
		es = append(es, e)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return es
}

// ClearLogFile empties the log file
func ClearLogFile() error {
	// Create the new file
	if _, err := os.Create(LogFile + ".new"); err != nil {
		Log.Error(err)
		return err
	}

	now := time.Now().Format(time.RFC3339)
	logTime, _ := time.Parse(time.RFC3339, now)
	newEntry := LogEntry{
		Level:   "warn",
		Message: "Truncated and cleared the log file!",
		Time:    logTime,
	}

	// Open and have the first entry informing user that the log has been recently emptied
	newLog, err := os.OpenFile(LogFile+".new", os.O_WRONLY, 0666)
	if err != nil {
		Log.Error(err)
		return err
	}
	jsonEntry, _ := json.Marshal(newEntry)
	if _, err := newLog.Write(jsonEntry); err != nil {
		Log.Error(err)
	}
	newLog.WriteString("\n")
	defer newLog.Close()

	// Rename the old error log, and update the name of the new one
	if err := os.Rename(LogFile, LogFile+".old"); err != nil {
		Log.Error("Could not rename and clear old log file!")
		return err
	}
	if err := os.Rename(LogFile+".new", LogFile); err != nil {
		Log.Error("Could not rename and clear old log file!")
		return err
	}

	if err := os.Remove(LogFile + ".old"); err != nil {
		Log.Error(err)
	}

	return nil

}

// ArchiveLogFile creates a backup of the current log in the logs directory, creates a new log file, and writes to the new log
func ArchiveLogFile() error {
	// Create the new file
	if _, err := os.Create(LogFile + ".new"); err != nil {
		Log.Error(err)
		return err
	}

	logTime := time.Now()
	newEntry := LogEntry{
		Level:   "warn",
		Message: "Archived the log file! Location: " + LogPath + strconv.Itoa(int(logTime.Unix())) + "-error.log",
		Time:    logTime,
	}

	// Open and have the first entry informing user that the log has been recently archived
	newLog, err := os.OpenFile(LogFile+".new", os.O_WRONLY, 0666)
	if err != nil {
		Log.Error(err)
		return err
	}
	jsonEntry, _ := json.Marshal(newEntry)
	if _, err = newLog.Write(jsonEntry); err != nil {
		Log.Error(err)
	}
	newLog.WriteString("\n")
	defer newLog.Close()

	// Rename the old error log, and update the name of the new one
	if err = os.Rename(LogFile, LogPath+strconv.Itoa(int(logTime.Unix()))+"-error.log"); err != nil {
		Log.Error("Could not archive log file!")
		return err
	}
	if err = os.Rename(LogFile+".new", LogFile); err != nil {
		Log.Error("Could not rename old log file!")
		return err
	}

	return err
}

// ToggleDebug toggles the debug level on or off depending on the current state
func ToggleDebug() {
	if Log.Level == logrus.DebugLevel {
		Log.WithFields(logrus.Fields{
			"Test": "Test",
		}).Info("Turned off debug logging...")
		Log.Level = logrus.InfoLevel
	} else {
		Log.WithFields(logrus.Fields{
			"Test": "Test",
		}).Info("Turned on debug logging...")
		Log.Level = logrus.DebugLevel
	}
}
