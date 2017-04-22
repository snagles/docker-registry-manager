package manager

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/Sirupsen/logrus"
)

const TestResourceDirectory = "/src/github.com/snagles/docker-registry-manager/tests"

// TestMain handles the setup and teardown of the test registry for use
func TestMain(m *testing.M) {
	cmd, err := SetupTestRegistry()
	if err != nil {
		logrus.Error("Unable to start the test registry!" + err.Error())
	}
	code := m.Run()
	err = TearDownTestRegistry(cmd)
	if err != nil {
		logrus.Error("Unable to teardown the registry!" + err.Error())
	}
	os.Exit(code)
}

// SetupTestRegistry by calling the registry serve command with test config
func SetupTestRegistry() (*exec.Cmd, error) {
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		goPath = defaultGOPATH()
	}

	// Check to see if the registry binary is installed
	if _, err := os.Stat(goPath + "/bin/registry"); os.IsNotExist(err) {
		logrus.WithFields(logrus.Fields{
			"Fix action":       "go get github.com/docker/distribution/cmd/registry",
			"More information": "https://github.com/docker/distribution/blob/master/BUILDING.md",
		}).Error("Registry binary is not installed in your GOPATH.")
	}

	os.Setenv("REGISTRY_STORAGE_FILESYSTEM_ROOTDIRECTORY", goPath+TestResourceDirectory+"/var/lib/registry")
	cmd := exec.Command(goPath+"/bin/registry", "serve", goPath+TestResourceDirectory+"/config-dev.yml")
	// start the command and dont wait
	err := cmd.Start()
	return cmd, err
}

// TearDownTestRegistry kills the serve command process started by the setup
func TearDownTestRegistry(cmd *exec.Cmd) error {
	err := cmd.Process.Kill()
	return err
}

// take from https://github.com/golang/go/blob/go1.8/src/go/build/build.go#L260-L277
func defaultGOPATH() string {
	env := "HOME"
	if runtime.GOOS == "windows" {
		env = "USERPROFILE"
	} else if runtime.GOOS == "plan9" {
		env = "home"
	}
	if home := os.Getenv(env); home != "" {
		def := filepath.Join(home, "go")
		if filepath.Clean(def) == filepath.Clean(runtime.GOROOT()) {
			// Don't set the default GOPATH to GOROOT,
			// as that will trigger warnings from the go tool.
			return ""
		}
		return def
	}
	return ""
}
