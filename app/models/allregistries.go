package manager

import (
	"net/url"
	"path"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// AllRegistries contains a list of added registries using their hostnames
// access granted via mutex locks/unlocks
var AllRegistries Registries

// avoidUpdateOnRegistries is a 'set' that stores names of registries that would
// not be updated in the current loop iteration that updates registries
var avoidUpdateOnRegistries = make(map[string]struct{})

func init() {
	AllRegistries.Registries = map[string]*Registry{}

	go func() {
		for {
			// This function can potentially hang the whole application
			// because it issues a write lock, that blocks all incomming http
			// calls until all registries would be updated.
			// Here is how this problem can be solved with minimum locking time
			// Step 1 -> lock application, collect names of registries that have
			// 			 to be updated, unlock application
			// Step 2 -> collect updated registry objects WITHOUT LOCKING.
			// 			 it is possible, because registry Update method has
			//           nonpointer receiver, so original registry object would
			// 			 not be updated
			// Step 3 -> lock application, replace old registries with updated
			//           ones

			// There are the following corner cases:
			// -> new registry was added between steps 1 and 3. This is not
			//    critical, because new registry does not have to be updated
			// -> registry was removed between steps 1 and 3. This can be
			//    critical, because it can resurrect deleted registries. To
			//    resolve this problem - additional check has to be performed
			//    before updating registries at step 3
			// -> registry was edited between steps 1 and 3. This is critical,
			//    because it can replace freshly edited registry with the old
			//    registry. To resolve this problem - module variable named
			//    avoidUpdateOnRegistries is introduced. This list
			//    stores names of registries that should not be updated at
			//    step 3
			// -> http call to refresh registry was made between steps 1 and 3.
			//    This can not cause any undefined behavior.

			avoidUpdateOnRegistries = make(map[string]struct{})
			registriesToUpdate := make(map[string]*Registry)

			// Step 1 -> Build list of names that have to be updated
			AllRegistries.Lock()
			for registryName, registry := range AllRegistries.Registries {
				if time.Now().UTC().Sub(registry.LastRefresh) >= registry.TTL {
					registriesToUpdate[registryName] = registry
				}
			}
			AllRegistries.Unlock()

			// Step 2 -> collect updated data in background
			for registryName, registry := range registriesToUpdate {
				updatedRegistry := registry.Update()
				registriesToUpdate[registryName] = &updatedRegistry
			}

			// Step 3 -> replace old registries with updated registries
			AllRegistries.Lock()
			for registryName, updatedRegistry := range registriesToUpdate {
				// check if registry was removed
				if _, registryExists := AllRegistries.Registries[registryName]; !registryExists {
					continue
				}

				// check if registry was edited
				if _, registryWasEdited := avoidUpdateOnRegistries[registryName]; registryWasEdited {
					continue
				}

				// registry can be safely updated
				AllRegistries.Registries[registryName] = updatedRegistry
			}
			AllRegistries.Unlock()

			time.Sleep(1 * time.Second)
		}
	}()
}

// Registries contains a map of all active registries identified by their name, locked when necessary
type Registries struct {
	Registries map[string]*Registry
	*viper.Viper
	sync.RWMutex
}

// Add adds a created registry to the map of AllRegistries using the corresponding locks
func (rs *Registries) Add(r *Registry) {
	AllRegistries.Lock()
	AllRegistries.Registries[r.Name] = r
	AllRegistries.Unlock()
	logrus.Infof("Added new registry: %s", r.Name)
}

// Remove removes a created registry from the map of AllRegistries using the corresponding locks
func (rs *Registries) Remove(r *Registry) {
	AllRegistries.Lock()
	delete(AllRegistries.Registries, r.Name)
	AllRegistries.Unlock()
	logrus.Infof("Removed registry: %s", r.Name)
}

// Edit updates the old registry with the new data
func (rs *Registries) Edit(new, old *Registry) {
	AllRegistries.Lock()
	delete(AllRegistries.Registries, old.Name)
	// copy the history
	new.History = old.History
	AllRegistries.Registries[new.Name] = new

	// make sure that this registry would avoid update loop if it is currently
	// running in the background
	avoidUpdateOnRegistries[new.Name] = struct{}{}

	AllRegistries.Unlock()
	logrus.Infof("Edited registry: %s", old.Name)
}

type registriesConfig struct {
	DisplayName          string
	URL                  string
	Port                 int
	Username             string
	Password             string
	SkipTLS              bool   `mapstructure:"skip-tls-validation" yaml:"skip-tls-validation"`
	RefreshRate          string `mapstructure:"refresh-rate" yaml:"refresh-rate"`
	DockerhubIntegration bool   `mapstructure:"dockerhub-integration" yaml:"dockerhub-integration"`
	ReadOnly             bool   `mapstructure:"read-only" yaml:"read-only"`
}

// LoadConfig adds the registries parsed from the passed yaml file
func (rs *Registries) LoadConfig(registriesFile string) {
	if rs.Viper == nil {
		rs.Viper = viper.New()
	}

	// If the registries path is not passed use the default project dir
	if registriesFile != "" {
		rs.AddConfigPath(path.Dir(registriesFile))
		base := path.Base(registriesFile)
		ext := path.Ext(registriesFile)
		rs.SetConfigName(base[0 : len(base)-len(ext)])
		logrus.Infof("Using registries located in %s with file name %s", path.Dir(registriesFile), base[0:len(base)-len(ext)])
	} else {
		rs.SetConfigName("registries")
		var root string
		_, run, _, ok := runtime.Caller(0)
		if ok {
			root = filepath.Dir(run)
			rs.AddConfigPath(root)
		} else {
			logrus.Fatalf("Failed to get runtime caller for parser")
		}
		logrus.Infof("Using registries located in %s with file name %s", root, "registries.yml")
	}

	config := make(map[string]map[string]registriesConfig)

	if err := rs.ReadInConfig(); err != nil {
		logrus.Fatalf("Failed to read in registries file: %s", err)
	}

	if err := rs.Unmarshal(&config); err != nil {
		logrus.Fatalf("Unable to unmarshal registries file: %s", err)
	}

	// overwrite the entries with the updated information
	for name, r := range config["registries"] {
		if r.URL != "" {
			url, err := url.Parse(r.URL)
			if err != nil {
				logrus.Fatalf("Failed to parse registry from the passed url (%s): %s", r.URL, err)
			}
			duration, err := time.ParseDuration(r.RefreshRate)
			if err != nil {
				logrus.Fatalf("Failed to add registry (%s), invalid duration: %s", r.URL, err)
			}
			new, err := NewRegistry(url.Scheme, url.Hostname(), name, r.DisplayName, r.Username, r.Password, r.Port, duration, r.SkipTLS, r.DockerhubIntegration, r.ReadOnly)
			if err != nil {
				logrus.Fatalf("Failed to create registry (%s): %s", r.URL, err)
			}
			AllRegistries.Add(new)
		}
	}
}

// WriteConfig builds the config and writes from the map of registries
func (rs *Registries) WriteConfig() error {
	config := make(map[string]registriesConfig)
	for name, r := range AllRegistries.Registries {
		config[name] = registriesConfig{
			DisplayName:          r.DisplayName,
			URL:                  r.URL,
			Port:                 r.Port,
			Username:             r.Username,
			Password:             r.Password,
			SkipTLS:              r.SkipTLS,
			ReadOnly:             r.ReadOnly,
			RefreshRate:          r.TTL.String(),
			DockerhubIntegration: r.DockerhubIntegration,
		}
	}

	rs.Set("registries", config)
	logrus.Info("Writing config with new/removed registries")
	return rs.Viper.WriteConfig()
}
