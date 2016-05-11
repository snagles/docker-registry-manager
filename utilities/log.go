package utils

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/rifflock/lfshook"
)

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

	// Create the log directory if needed
	if _, err := os.Stat("./logs"); os.IsNotExist(err) {
		if err := os.Mkdir("logs", os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}

	// Create the log file if needed
	if _, err := os.Stat("./logs/error.log"); os.IsNotExist(err) {
		if _, err = os.Create("./logs/error.log"); err != nil {
			log.Fatal(err)
		}
	}

	Log.Formatter = &logrus.JSONFormatter{}
	Log.Hooks.Add(lfshook.NewHook(lfshook.PathMap{
		logrus.ErrorLevel: "./logs/error.log",
		logrus.InfoLevel:  "./logs/error.log",
		logrus.WarnLevel:  "./logs/error.log",
		logrus.FatalLevel: "./logs/error.log",
		logrus.PanicLevel: "./logs/error.log",
		logrus.DebugLevel: "./logs/error.log",
	}))
}

// ParseLogFile parses the locally stored flat log file that was logged to by logrus
func ParseLogFile() []LogEntry {
	file, err := os.Open("logs/error.log")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	es := []LogEntry{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		e := LogEntry{}
		err := json.Unmarshal([]byte(scanner.Text()), &e)
		if err != nil {
			Log.Error(err)
		}
		es = append(es, e)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return es
}

func ClearLogFile() error {
	// Create the new file
	if _, err := os.Create("./logs/error.log.new"); err != nil {
		Log.Error(err)
		return err
	}
	newEntry := LogEntry{
		Level:   "warn",
		Message: "Truncated and cleared the log file!",
		Time:    time.Now(),
	}

	// Open and have the first entry informing user that the log has been recently emptied
	newLog, err := os.Open("logs/error.log.new")
	if err != nil {
		Log.Error(err)
		return err
	}
	jsonEntry, _ := json.Marshal(newEntry)
	newLog.Write(jsonEntry)
	defer newLog.Close()

	// Rename the old error log, and update the name of the new one
	if err := os.Rename("logs/error.log", "logs/error.log.old"); err != nil {
		Log.Error("Could not rename and clear old log file!")
		return err
	}
	if err := os.Rename("logs/error.log.new", "logs/error.log"); err != nil {
		Log.Error("Could not rename and clear old log file!")
		return err
	}

	return nil

}
