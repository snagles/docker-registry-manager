package registry

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	client "github.com/snagles/docker-registry-client/registry"
	"github.com/snagles/docker-registry-manager/app/models"
)

// RegistriesController extends the beego.Controller type
type RegistriesController struct {
	beego.Controller
}

// Get returns the template for the registries page
func (c *RegistriesController) Get() {
	manager.AllRegistries.RLock()
	c.Data["registries"] = manager.AllRegistries.Registries
	manager.AllRegistries.RUnlock()
	// Index template
	c.TplName = "registries.tpl"
}

// GetRegistryCount returns the number of currently added registries
func (c *RegistriesController) GetRegistryCount() {
	manager.AllRegistries.RLock()
	c.Data["registries"] = manager.AllRegistries.Registries

	registryCount := struct {
		Count int
	}{
		len(manager.AllRegistries.Registries),
	}
	c.Data["json"] = &registryCount
	manager.AllRegistries.RUnlock()
	c.ServeJSON()
}

// AddRegistry adds a registry to the active registry list from a form
func (c *RegistriesController) AddRegistry() {
	// Registry contains all identifying information for communicating with a registry
	scheme, host, name, displayName, port, skipTLS, dockerhubIntegration, readOnly, err := c.sanitizeForm()
	if err != nil {
		c.CustomAbort(404, err.Error())
	}

	interval, err := c.GetInt("interval", 60)
	if err != nil {
		c.CustomAbort(404, err.Error())
	}

	ttl := time.Duration(interval) * time.Second

	r, err := manager.NewRegistry(scheme, host, name, displayName, "", "", port, ttl, skipTLS, dockerhubIntegration, readOnly)
	if err != nil {
		c.CustomAbort(404, err.Error())
	}
	manager.AllRegistries.Add(r)
	manager.AllRegistries.WriteConfig()
	c.Ctx.Redirect(302, "/registries")
}

func (c *RegistriesController) EditRegistry() {
	// Registry contains all identifying information for communicating with a registry
	scheme, host, name, displayName, port, skipTLS, dockerhubIntegration, readOnly, err := c.sanitizeForm()
	if err != nil {
		c.CustomAbort(404, err.Error())
	}

	interval, err := c.GetInt("interval", 60)
	if err != nil {
		c.CustomAbort(404, err.Error())
	}

	new, err := manager.NewRegistry(scheme, host, name, displayName, "", "", port, time.Duration(interval)*time.Second, skipTLS, dockerhubIntegration, readOnly)
	if err != nil {
		c.CustomAbort(404, err.Error())
	}

	registryName := c.Ctx.Input.Param(":registryName")
	if old, ok := manager.AllRegistries.Registries[registryName]; ok {
		manager.AllRegistries.Edit(new, old)
	}

	manager.AllRegistries.WriteConfig()
	c.Ctx.Redirect(302, "/registries")
}

// RegistryStatus responds with JSON containing the status of the registry
func (c *RegistriesController) RegistryStatus() {
	// Define the response
	var res struct {
		Error       string `json:"error, omitempty"`
		IsAvailable bool   `json:"is_available"`
	}
	var err error
	whenErr := func(err error) {
		res.Error = err.Error()
		res.IsAvailable = false
		c.Data["json"] = &res
		c.ServeJSON()
	}

	scheme, host, _, _, port, skipTLS, _, _, err := c.sanitizeForm()
	if err != nil {
		whenErr(err)
		return
	}

	url := fmt.Sprintf("%s://%s:%v", scheme, host, port)
	var temp *client.Registry
	if skipTLS {
		temp, err = client.NewInsecure(url, "", "")
	} else {
		temp, err = client.New(url, "", "")
	}
	if err != nil {
		whenErr(err)
		return
	}
	err = temp.Ping()
	if err != nil {
		whenErr(err)
		return
	}

	res.Error = ""
	res.IsAvailable = true
	c.Data["json"] = &res
	c.ServeJSON()
}

// Refresh refreshes the passed registry
func (c *RegistriesController) Refresh() {
	registryName := c.Ctx.Input.Param(":registryName")

	manager.AllRegistries.RLock()
	r, ok := manager.AllRegistries.Registries[registryName]
	manager.AllRegistries.RUnlock()

	if ok {
		updatedRegistry := r.Update()
		// Registry is being updated, so write lock has to be used instead of read
		manager.AllRegistries.Lock()
		manager.AllRegistries.Registries[registryName] = &updatedRegistry
		manager.AllRegistries.Unlock()
	}
	// Index template
	c.CustomAbort(200, "Refreshed registry")
}

func (c *RegistriesController) sanitizeForm() (scheme, host string, name string, displayName string, port int, skipTLS bool, dockerhubIntegration bool, readOnly bool, err error) {
	host = c.GetString("host")
	port, err = c.GetInt("port", 5000)
	name = c.GetString("name", fmt.Sprintf("%s:%v", host, port))
	displayName = c.GetString("displayName")
	scheme = c.GetString("scheme", "https")
	skipTLSOn := c.GetString("skip-tls-validation", "off")
	if skipTLSOn == "on" {
		skipTLS = true
	}

	dockerhubIntegrationOn := c.GetString("dockerhub-integration", "off")
	if dockerhubIntegrationOn == "on" {
		dockerhubIntegration = true
	}

	readOnlyOn := c.GetString("read-only", "off")
	if readOnlyOn == "on" {
		readOnly = true
	}

	switch {
	case scheme == "":
		err = errors.New("Invalid scheme: " + scheme)
		return
	case host == "":
		err = errors.New("Invalid host: " + host)
		return
	case port == 0:
		err = errors.New("Invalid port: " + strconv.Itoa(port))
		return
	}
	return
}
