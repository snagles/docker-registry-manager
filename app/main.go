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
			Usage:       "tls certificate path path and name `/app/key.key`",
			EnvVar:      "MANAGER_KEY",
			Destination: &keyPath,
		},
		cli.StringFlag{
			Name:        "tls-certificate, cert",
			Usage:       "tls certificate path path and name `/app/certificate.crt`",
			EnvVar:      "MANAGER_CERTIFICATE",
			Destination: &certPath,
		},
	}

	app.Action = func(ctx *cli.Context) {
		c, err := parseRegistries(registriesFile)
		if err != nil {
			logrus.Fatal(err)
		}

		err = setlevel(logLevel)
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

		// set http port
		if enableHTTPS {
			beego.BConfig.Listen.HTTPSPort = appPort
			// make sure we have both key and cert
			if beego.BConfig.Listen.HTTPSKeyFile == "" {
				logrus.Fatal("HTTPS enabled, but no key file provided")
			} else {
				beego.BConfig.Listen.HTTPSKeyFile = keyPath
			}
			if beego.BConfig.Listen.HTTPSKeyFile == "" {
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

type registries struct {
	Registries map[string]struct {
		URL         string
		Username    string
		Password    string
		SkipTLS     bool   `mapstructure:"skip-tls-validation"`
		RefreshRate string `mapstructure:"refresh-rate"`
	} `mapstructure:"registries"`
}

func parseRegistries(registriesFile string) (*registries, error) {
	v := viper.New()

	// If the registries path is not passed use the default project dir
	if registriesFile != "" {
		v.AddConfigPath(path.Dir(registriesFile))
		base := path.Base(registriesFile)
		ext := path.Ext(registriesFile)
		v.SetConfigName(base[0 : len(base)-len(ext)])
		logrus.Infof("Using registries located in %s with file name %s", path.Dir(registriesFile), base[0:len(base)-len(ext)])
	} else {
		v.SetConfigName("registries")
		var root string
		_, r, _, ok := runtime.Caller(0)
		if ok {
			root = filepath.Dir(r)
			v.AddConfigPath(root)
		} else {
			logrus.Fatalf("Failed to get runtime caller for parser")
		}
		logrus.Infof("Using registries located in %s with file name %s", root, "registries.yml")
	}

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("Failed to read in registries file: %s", err)
	}

	c := registries{}
	if err := v.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("Unable to unmarshal registries file: %s", err)
	}
	return &c, nil
}
