package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/astaxie/beego"
	"github.com/snagles/docker-registry-manager/app/models/manager"
	_ "github.com/snagles/docker-registry-manager/app/routers"
	funcs "github.com/snagles/docker-registry-manager/app/templatefuncs"
	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "Docker Registry Manager"
	app.Usage = "Connect to, view, and manage multiple private Docker registries"
	app.Version = "2.0.0"
	var logLevel string
	var refreshRate string

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
			Name:   "registries, r",
			Usage:  "comma separated list of registry url's to connect to `http://url:5000,https://url:6000`",
			EnvVar: "MANAGER_REGISTRIES",
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
	}

	app.Action = func(c *cli.Context) {
		setlevel(logLevel)

		// set default beego config
		beego.BConfig.AppName = "docker-registry-manager"
		beego.BConfig.RunMode = "dev"
		beego.BConfig.Listen.EnableAdmin = true
		beego.BConfig.CopyRequestBody = true

		// add template functions
		beego.AddFuncMap("shortenDigest", funcs.DigestShortener)
		beego.AddFuncMap("bytefmt", funcs.ByteFmt)
		beego.AddFuncMap("bytefmtdiff", funcs.ByteDiffFmt)
		beego.AddFuncMap("timeAgo", funcs.TimeAgo)
		beego.AddFuncMap("oneIndex", func(i int) int { return i + 1 })

		registries := strings.Split(c.String("registries"), ",")
		for _, registry := range registries {
			if registry != "" {
				url, err := url.Parse(registry)
				if err != nil {
					fmt.Print("Failed to parse registry from the passed url.\n\n")
					cli.ShowAppHelp(c)
					return
				}
				port, err := strconv.Atoi(url.Port())
				if err != nil || port == 0 {
					fmt.Print("Failed to add registry, invalid port!\n\n")
					cli.ShowAppHelp(c)
					return
				}
				duration, err := time.ParseDuration(refreshRate)
				if err != nil {
					fmt.Print("Failed to add registry, invalid duration!\n\n")
					cli.ShowAppHelp(c)
					return
				}
				manager.AddRegistry(url.Scheme, url.Hostname(), port, duration)
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
