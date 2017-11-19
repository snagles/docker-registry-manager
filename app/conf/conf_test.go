package conf

import (
	"os"
	"testing"
)

// TestGoPath tests that the GOPATH is set correctly
func TestGoPath(t *testing.T) {
	if GOPATH == "" {
		t.Error("We should be able to retrieve the current gopath")
	}
}

// TestLogFileCreation tests that the logfile is successfully created
func TestLogFileCreation(t *testing.T) {
	if LogDir == "" || LogFile == "" {
		t.Error("The log file and directory should be set correctly")
	}

	_, err := os.Stat(LogFile)
	if err != nil {
		t.Error("The log file and directory should be created successfully")
	}
}
