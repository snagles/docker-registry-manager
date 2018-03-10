package main

import (
	"fmt"
	"os"

	"github.com/astaxie/beego"
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/snagles/docker-registry-manager/app/models"
	_ "github.com/snagles/docker-registry-manager/app/routers"
	"github.com/urfave/cli"
)

const (
	appVersion = "3.0.0"
)

func main() {
	app := cli.NewApp()
	app.Name = "Docker Registry Manager"
	app.Usage = "Connect to, view, and manage multiple private Docker registries"

	var registriesFile, logLevel, keyPath, certPath string
	var enableHTTPS bool
	var appPort int
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:        "port, p",
			Usage:       "port to use for the registry manager `port`",
			Value:       8080,
			Destination: &appPort,
			EnvVar:      "MANAGER_PORT",
		},
		cli.StringFlag{
			Name:        "registries, r",
			Usage:       "file location of the registries.yml `/app/registries.yml`",
			EnvVar:      "MANAGER_REGISTRIES",
			Destination: &registriesFile,
		},
		cli.StringFlag{
			Name:        "log-level, l",
			Usage:       "log-level `warn`",
			Value:       "info",
			EnvVar:      "MANAGER_LOG_LEVEL",
			Destination: &logLevel,
		},

		// Beego HTTPS options
		cli.BoolFlag{
			Name:        "enable-https, e",
			Usage:       "enable https `true or false`",
			EnvVar:      "MANAGER_ENABLE_HTTPS",
			Destination: &enableHTTPS,
		},
		cli.StringFlag{
			Name:        "tls-key, k",
			Usage:       "tls certificate path and name `/app/key.key`",
			EnvVar:      "MANAGER_KEY",
			Destination: &keyPath,
		},
		cli.StringFlag{
			Name:        "tls-certificate, cert",
			Usage:       "tls certificate path and name `/app/certificate.crt`",
			EnvVar:      "MANAGER_CERTIFICATE",
			Destination: &certPath,
		},
	}

	app.Action = func(ctx *cli.Context) {

		manager.AllRegistries.LoadConfig(registriesFile)

		err := setlevel(logLevel)
		if err != nil {
			logrus.Fatal(err)
		}

		// Beego configuration
		beego.BConfig.AppName = "docker-registry-manager"
		beego.BConfig.RunMode = "dev"
		beego.BConfig.Listen.EnableAdmin = true
		beego.BConfig.CopyRequestBody = true
		beego.BConfig.WebConfig.ViewsPath = "views"

		// set http port
		if enableHTTPS {
			beego.BConfig.Listen.HTTPSPort = appPort
			// make sure we have both key and cert
			if keyPath == "" {
				logrus.Fatal("HTTPS enabled, but no key file provided")
			} else {
				beego.BConfig.Listen.HTTPSKeyFile = keyPath
			}
			if certPath == "" {
				logrus.Fatal("HTTPS enabled, but no certificate file provided")
			} else {
				beego.BConfig.Listen.HTTPSCertFile = certPath
			}
			// if we're not using https just use standard http
		} else {
			beego.BConfig.Listen.HTTPPort = appPort
		}

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

func addRegistries(c *registries) {
	for name, r := range c.Registries {
		if r.URL != "" {
			url, err := url.Parse(r.URL)
			if err != nil {
				logrus.Fatalf("Failed to parse registry from the passed url (%s): %s", r.URL, err)
			}
			duration, err := time.ParseDuration(r.RefreshRate)
			if err != nil {
				logrus.Fatalf("Failed to add registry (%s), invalid duration: %s", r.URL, err)
			}
			if r.Password != "" && r.Username != "" {
				if _, err := manager.AddRegistry(url.Scheme, url.Hostname(), name, r.Username, r.Password, r.Port, duration, r.SkipTLS, r.DockerhubIntegration); err != nil {
					logrus.Fatalf("Failed to add registry (%s): %s", r.URL, err)
				}
			} else {
				if _, err := manager.AddRegistry(url.Scheme, url.Hostname(), name, "", "", r.Port, duration, r.SkipTLS, r.DockerhubIntegration); err != nil {
					logrus.Fatalf("Failed to add registry (%s): %s", r.URL, err)
				}
			}
		}
	}
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

