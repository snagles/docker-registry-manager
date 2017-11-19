package main

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/sirupsen/logrus"
	"github.com/snagles/docker-registry-manager/app/conf"
	"github.com/snagles/docker-registry-manager/app/models"
	_ "github.com/snagles/docker-registry-manager/app/routers"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "Docker Registry Manager"
	app.Usage = "Connect to, view, and manage multiple private Docker registries"
	app.Version = "2.0.1"
	var logLevel string
	var refreshRate string
	var skipTLS bool
	var passedRegistries string
	var configFile string

	cli.AppHelpTemplate = fmt.Sprintf(`%s
WEBSITE:
  https://github.com/snagles/docker-registry-manager
	`, cli.AppHelpTemplate)

	app.Authors = []cli.Author{
		cli.Author{
			Name: "Stefan Naglee",
		},
	}

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:        "port, p",
			Usage:       "port to use for the registry manager `port`",
			Value:       8080,
			Destination: &beego.BConfig.Listen.HTTPPort,
			EnvVar:      "MANAGER_PORT",
		},
		cli.StringFlag{
			Name:        "registries, r",
			Usage:       "comma separated list of registry url's to connect to `http://url:5000,https://url:6000`",
			EnvVar:      "MANAGER_REGISTRIES",
			Destination: &passedRegistries,
		},
		cli.StringFlag{
			Name:        "registries-config, rc",
			Usage:       "file location of the config.yml `~/config.yml`",
			EnvVar:      "MANAGER_REGISTRIES_CONFIG",
			Destination: &configFile,
		},
		cli.StringFlag{
			Name:        "log, l",
			Usage:       "log level `level`",
			Value:       "info",
			EnvVar:      "MANAGER_LOG_LEVEL",
			Destination: &logLevel,
		},
		cli.StringFlag{
			Name:        "ttl, t",
			Usage:       "ttl refresh rate `h,m,s,ms`",
			Value:       "30s",
			EnvVar:      "MANAGER_REFRESH_RATE",
			Destination: &refreshRate,
		},
		cli.BoolFlag{
			Name:        "skip-tls",
			Usage:       "skip-tls",
			EnvVar:      "MANAGER_SKIP_TLS",
			Destination: &skipTLS,
		},
	}

	app.Action = func(c *cli.Context) {
		setlevel(logLevel)

		// set default beego config
		beego.BConfig.AppName = "docker-registry-manager"
		beego.BConfig.RunMode = "dev"
		beego.BConfig.Listen.EnableAdmin = true
		beego.BConfig.CopyRequestBody = true

		// add template functions
		beego.AddFuncMap("shortenDigest", DigestShortener)
		beego.AddFuncMap("statToSeconds", StatToSeconds)
		beego.AddFuncMap("bytefmt", ByteFmt)
		beego.AddFuncMap("bytefmtdiff", ByteDiffFmt)
		beego.AddFuncMap("timeAgo", TimeAgo)
		beego.AddFuncMap("oneIndex", func(i int) int { return i + 1 })

		confs := make(map[string]config)
		if configFile != "" {
			confs = parseConfigs(configFile)
		} else {
			registries := strings.Split(passedRegistries, ",")
			for _, registry := range registries {
				confs[registry] = config{URL: registry}
			}
		}

		for _, conf := range confs {
			if conf.URL != "" {
				url, err := url.Parse(conf.URL)
				if err != nil {
					fmt.Printf("Failed to parse registry from the passed url: %s \n\n", err)
					cli.ShowAppHelp(c)
					return
				}
				port, err := strconv.Atoi(url.Port())
				if err != nil || port == 0 {
					fmt.Printf("Failed to add registry, invalid port: %s\n\n", err)
					cli.ShowAppHelp(c)
					return
				}
				duration, err := time.ParseDuration(refreshRate)
				if err != nil {
					fmt.Printf("Failed to add registry, invalid duration: %s \n\n", err)
					cli.ShowAppHelp(c)
					return
				}
				if conf.Password != "" && conf.Username != "" {
					_, err := manager.AddRegistry(url.Scheme, url.Hostname(), conf.Username, conf.Password, port, duration, skipTLS)
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
				} else {
					_, err := manager.AddRegistry(url.Scheme, url.Hostname(), "", "", port, duration, skipTLS)
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
				}
			}
		}
		beego.Run()
	}
	app.Run(os.Args)
}

func setlevel(level string) {
	switch {
	case level == "panic":
		logrus.SetLevel(logrus.PanicLevel)
	case level == "fatal":
		logrus.SetLevel(logrus.FatalLevel)
	case level == "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case level == "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case level == "info":
		logrus.SetLevel(logrus.InfoLevel)
	case level == "debug":
		logrus.SetLevel(logrus.DebugLevel)
	}
}

type config struct {
	URL      string
	Username string
	Password string
}

// parseConfigs parses the yaml databases config and add them to the Configs map
func parseConfigs(configPath string) map[string]config {
	var c = viper.New()

	// If the config path is not passed use the default project dir
	if configPath != "" {
		c.AddConfigPath(path.Dir(configPath))
		base := path.Base(configPath)
		ext := path.Ext(configPath)
		c.SetConfigName(base[0 : len(base)-len(ext)])
	} else {
		// use the default tree
		c.AddConfigPath(conf.GOPATH + "/src/github.com/snagles/docker-registry-manager")
		c.SetConfigName("config")
	}

	if err := c.ReadInConfig(); err != nil {
		logrus.Fatalf("Fatal error config file: %v", err)
	}

	var confs map[string]config
	if err := c.Unmarshal(&confs); err != nil {
		logrus.Fatalf("unable to decode into struct, %v", err)
	}
	return confs
}
