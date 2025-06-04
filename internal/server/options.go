package server

import (
	"os"

	goflags "github.com/jessevdk/go-flags"
)

type Options struct {
	ApiPort       string `long:"api-port" description:"The port to use for api calls" default:"8081"`
	SitePort      string `short:"p" long:"port" description:"The port to use for the dashboard" default:"8080"`
	Debug         bool   `long:"debug" description:"Print debug messages"`
	Demo          bool   `long:"run-demo" description:"Run a demo, will use temporary folder"`
	DefaultStatic bool   `long:"default-static" description:"Use default static content"`
	DataPath      string `short:"d" long:"data-path" description:"Path to store data files (default: $HOME/.alive/data)"`
	StaticPath    string `long:"static-path" description:"Path to store static files (default: $HOME/.alive/static)"`
	ParentUrl     string `long:"parent-url" description:"Url for a parent dashboard, if set enables updating a parent dashboard with the overal status of this dashboard"`
	ParentBoxID   string `long:"parent-id" description:"Box id to use when updating status on a parent dashboard"`
	ParentBoxSize string `long:"parent-size" description:"Box size to use when updating status on a parent dashboard (default: large)" default:"large"`
}

var options Options

func processOptions() {
	goflagParser := goflags.NewParser(&options, goflags.Default)

	if _, err := goflagParser.Parse(); err != nil {
		if flagsErr, ok := err.(*goflags.Error); ok && flagsErr.Type == goflags.ErrHelp {
			os.Exit(0)
		} else {
			logger.Panic(err.Error())
		}
	}
}
