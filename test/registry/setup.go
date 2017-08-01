package registry

import (
	"fmt"
	// Used to add profiling for the registry if debug is enabled
	_ "net/http/pprof"
	"os/exec"

	"github.com/Sirupsen/logrus"
	"github.com/docker/distribution/registry"
	// Setup auth and storage drivers
	_ "github.com/docker/distribution/registry/auth/htpasswd"
	_ "github.com/docker/distribution/registry/auth/silly"
	_ "github.com/docker/distribution/registry/auth/token"
	_ "github.com/docker/distribution/registry/proxy"
	_ "github.com/docker/distribution/registry/storage/driver/filesystem"
	_ "github.com/docker/distribution/registry/storage/driver/gcs"
	_ "github.com/docker/distribution/registry/storage/driver/inmemory"
)

var stop chan struct{}

// Repositories contains a list of test repos
var Repositories = map[string][]string{
	"golang": []string{"1.7", "alpine"},
}

// Start creates and executes a test registry running with the
// specified parameters passing in the config
func Start(config string) {
	// Use the same configuration and setup run when using the registry binary
	// from the command line
	cmd := registry.ServeCmd
	stop = make(chan struct{})
	go func() {
		go cmd.Run(cmd, []string{config})
		for {
			select {
			case <-stop:
				break
			}
		}
	}()
}

// Stop signals the currently running test registry to quit
func Stop() {
	stop <- struct{}{}
	close(stop)
}

// Seed adds a subset of repositories for testing
func Seed(url string) error {

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
