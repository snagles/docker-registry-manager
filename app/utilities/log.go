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

//{"level":"error","msg":"Test","time":"2016-04-10T12:05:30-04:00"}
type LogEntry struct {
	Level string
	Msg   string
	Time  time.Time
}

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
