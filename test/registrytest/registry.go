package registrytest

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	// Used to add profiling for the registry if debug is enabled
	_ "net/http/pprof"
	"os/exec"

	"github.com/sirupsen/logrus"
	// Setup auth and storage drivers
	"github.com/docker/distribution/configuration"
	"github.com/docker/distribution/registry"
	_ "github.com/docker/distribution/registry/auth/htpasswd"
	_ "github.com/docker/distribution/registry/auth/silly"
	_ "github.com/docker/distribution/registry/auth/token"
	_ "github.com/docker/distribution/registry/proxy"
	_ "github.com/docker/distribution/registry/storage/driver/filesystem"
	_ "github.com/docker/distribution/registry/storage/driver/gcs"
	_ "github.com/docker/distribution/registry/storage/driver/inmemory"
)

// Repositories contains a list of test repos
var Repositories = map[string][]string{
	"golang": []string{"1.7", "alpine"},
}

// Start creates and executes a test registry running with the
// specified parameters passed in the config
func New(configType string) *registry.Registry {
	GOPATH := os.Getenv("GOPATH")
	if GOPATH == "" {
		GOPATH = defaultGOPATH()
	}
	r, err := os.Open(GOPATH + "/src/github.com/snagles/docker-registry-manager/test/registrytest/configs/" + configType)
	if err != nil {
		logrus.Fatalf("Failed to open %v configuration file: %v", configType, err)
	}
	c, err := configuration.Parse(r)
	if err != nil {
		logrus.Fatalf("Failed to parse %v configuration file: %v", c, err)
	}

	registry, err := registry.NewRegistry(context.Background(), c)
	if err != nil {
		logrus.Fatal(err)
	}

	return registry
}

func Start(r *registry.Registry) error {
	var err error
	go func(err error) {
		if err = r.ListenAndServe(); err != nil {
			logrus.Error(err)
		}
	}(err)
	return err
}

// Stop shutsdown the passed registry
func Stop(r *registry.Registry) error {
	return r.Server.Shutdown(context.Background())
}

// Seed adds a subset of repositories for testing
func Seed(r *registry.Registry) error {
	var url string
	if r.Config.HTTP.Host+":"+r.Config.HTTP.Addr == "::5000" {
		url = "localhost:5000"
	}
	for repo, tags := range Repositories {
		for _, tag := range tags {

			hubTag := fmt.Sprintf("%s:%s", repo, tag)
			registryTag := fmt.Sprintf("%s/%s:%s", url, repo, tag)

			//docker pull ubuntu:latest
			if err := exec.Command("docker", "pull", hubTag).Run(); err != nil {
				logrus.Errorf("Failed to pull: %s error: %s", hubTag, err.Error())
				return err
			}

			//docker tag ubuntu:latest localhost:5000/ubuntu:latest
			if err := exec.Command("docker", "tag", hubTag, registryTag).Run(); err != nil {
				logrus.Errorf("Failed to tag: %s %s error: %s", hubTag, registryTag, err.Error())
				return err
			}

			//docker push localhost:5000/ubuntu:latest
			if err := exec.Command("docker", "push", registryTag).Run(); err != nil {
				logrus.Errorf("Failed to push: %s error: %s", registryTag, err.Error())
				return err
			}
		}
	}
	return nil
}

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
