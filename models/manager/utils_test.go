package manager

import (
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/docker/distribution/configuration"
	"github.com/docker/distribution/context"
	"github.com/docker/distribution/registry"
	_ "github.com/docker/distribution/registry/storage/driver/filesystem"
	"github.com/docker/distribution/version"
)

const TestResourceDirectory = "/src/github.com/snagles/docker-registry-manager/tests"

// TestMain handles the setup and teardown of the test registry for use
func TestMain(m *testing.M) {
	go SetupTestRegistry()
	var ready bool
	for ready == false {
		resp, _ := http.Get("http://localhost:5010/v2/")
		if resp != nil {
			ready = true
			defer resp.Body.Close()
		}
	}
	code := m.Run()
	os.Exit(code)
}

// SetupTestRegistry by calling the registry serve command with test config
func SetupTestRegistry() {
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		goPath = defaultGOPATH()
	}

	// Use the test config file to create and run a new registry
	ctx := context.WithVersion(context.Background(), version.Version)
	fp, err := os.Open(goPath + TestResourceDirectory + "/config-dev.yml")
	if err != nil {
		logrus.Fatalln(err)
	}
	devConfig, err := configuration.Parse(fp)
	if err != nil {
		logrus.Fatalln(err)
	}
	// overwrite the filesystem location to use the project path
	devConfig.Storage["filesystem"]["rootdirectory"] = goPath + TestResourceDirectory + "/var/lib/registry"

	// Create the new registry
	registry, err := registry.NewRegistry(ctx, devConfig)
	if err != nil {
		logrus.Fatalln(err)
	}

	// Start the registry
	err = registry.ListenAndServe()
	if err != nil {
		logrus.Fatalln(err)
	}
}

// taken from https://github.com/golang/go/blob/go1.8/src/go/build/build.go#L260-L277
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
