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

	c.Data["registries"] = manager.AllRegistries.Registries

	// Index template
	c.TplName = "registries.tpl"
}

func (c *RegistriesController) GetRegistryCount() {
	c.Data["registries"] = manager.AllRegistries.Registries

	registryCount := struct {
		Count int
	}{
		len(manager.AllRegistries.Registries),
	}
	c.Data["json"] = &registryCount
	c.ServeJSON()
}

// AddRegistry adds a registry to the active registry list from a form
func (c *RegistriesController) AddRegistry() {
	// Registry contains all identifying information for communicating with a registry

	scheme, host, port, skipTLS, err := c.sanitizeForm()
	if err != nil {
		c.CustomAbort(404, err.Error())
	}

	interval, err := c.GetInt("interval", 10)
	if err != nil {
		c.CustomAbort(404, err.Error())
	}

	ttl := time.Duration(interval) * time.Second

	_, err = manager.AddRegistry(scheme, host, "", "", port, ttl, skipTLS)
	if err != nil {
		c.CustomAbort(404, err.Error())
	}
	c.Ctx.Redirect(302, "/registries")
}

// TestRegistryStatus responds with JSON containing the status of the registry
func (c *RegistriesController) TestRegistryStatus() {
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

	scheme, host, port, skipTLS, err := c.sanitizeForm()
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

func (c *RegistriesController) sanitizeForm() (scheme, host string, port int, skipTLS bool, err error) {
	host = c.GetString("host")
	port, err = c.GetInt("port", 5000)
	scheme = c.GetString("scheme", "https")
	skipTLSOn := c.GetString("skip-tls", "off")
	if skipTLSOn == "on" {
		skipTLS = true
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
