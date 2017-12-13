package main

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/sirupsen/logrus"
	"github.com/snagles/docker-registry-manager/app/models"
	_ "github.com/snagles/docker-registry-manager/app/routers"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

const (
	appVersion = "2.0.1"
)

func main() {
	app := cli.NewApp()
	app.Name = "Docker Registry Manager"
	app.Usage = "Connect to, view, and manage multiple private Docker registries"

	var configPath string
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config-path, c",
			Usage:       "config path and name `/opt/docker-registry-manager/config.yml`",
			EnvVar:      "REGISTRY_CONFIG",
			Destination: &configPath,
		},
	}

	app.Action = func(ctx *cli.Context) {
		c, err := parseConfig(configPath)
		if err != nil {
			logrus.Fatal(err)
		}

		err = setlevel(c.App.LogLevel)
		if err != nil {
			logrus.Fatal(err)
		}

		for _, r := range c.Registries {
			if r.URL != "" {
				url, err := url.Parse(r.URL)
				if err != nil {
					logrus.Fatalf("Failed to parse registry from the passed url (%s): %s", r.URL, err)
				}
				port, err := strconv.Atoi(url.Port())
				if err != nil || port == 0 {
					logrus.Fatalf("Failed to add registry (%s), invalid port: %s", r.URL, err)
				}
				duration, err := time.ParseDuration(r.RefreshRate)
				if err != nil {
					logrus.Fatalf("Failed to add registry (%s), invalid duration: %s", r.URL, err)
				}
				if r.Password != "" && r.Username != "" {
					if _, err := manager.AddRegistry(url.Scheme, url.Hostname(), r.Username, r.Password, port, duration, r.SkipTLS); err != nil {
						logrus.Fatalf("Failed to add registry (%s): %s", r.URL, err)
					}
				} else {
					if _, err := manager.AddRegistry(url.Scheme, url.Hostname(), "", "", port, duration, r.SkipTLS); err != nil {
						logrus.Fatalf("Failed to add registry (%s): %s", r.URL, err)
					}
				}
			}
		}

		// Beego configuration
		beego.BConfig.AppName = "docker-registry-manager"
		beego.BConfig.RunMode = "dev"
		beego.BConfig.Listen.EnableAdmin = true
		beego.BConfig.CopyRequestBody = true
		beego.BConfig.WebConfig.ViewsPath = "views"

		// add template functions
		beego.AddFuncMap("shortenDigest", DigestShortener)
		beego.AddFuncMap("statToSeconds", StatToSeconds)
		beego.AddFuncMap("bytefmt", ByteFmt)
		beego.AddFuncMap("bytefmtdiff", ByteDiffFmt)
		beego.AddFuncMap("timeAgo", TimeAgo)
		beego.AddFuncMap("oneIndex", func(i int) int { return i + 1 })
		beego.Run()
	}
	app.Run(os.Args)
}

func setlevel(level string) error {
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
	default:
		return fmt.Errorf("Unrecognized log level: %s", level)
	}
	return nil
}

type config struct {
	App struct {
		LogLevel string `mapstructure:"log-level"`
		Port     int
	}
	Registries map[string]struct {
		URL         string
		Username    string
		Password    string
		SkipTLS     bool   `mapstructure:"skip-tls-validation"`
		RefreshRate string `mapstructure:"refresh-rate"`
	} `mapstructure:"registries"`
}

func parseConfig(configPath string) (*config, error) {
	v := viper.New()

	// If the config path is not passed use the default project dir
	if configPath != "" {
		v.AddConfigPath(path.Dir(configPath))
		base := path.Base(configPath)
		ext := path.Ext(configPath)
		v.SetConfigName(base[0 : len(base)-len(ext)])
		logrus.Infof("Using config located in %s with config name %s", path.Dir(configPath), base[0:len(base)-len(ext)])
	} else {
		v.SetConfigName("config")
		var root string
		_, r, _, ok := runtime.Caller(0)
		if ok {
			root = filepath.Dir(r)
			v.AddConfigPath(root)
		} else {
			logrus.Fatalf("Failed to get runtime caller for parser")
		}
		logrus.Infof("Using config located in %s with config name %s", root, "config.yml")
	}

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("Failed to read in config file: %s", err)
	}

	c := config{}
	if err := v.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("Unable to unmarshal config file: %s", err)
	}
	return &c, nil
}
