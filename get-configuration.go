package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Config struct contains configuration to be used throughout the program
type Config struct {
	apiPort          string
	baseDir          string
	dataFile         string
	sitePort         string
	staticFilePath   string
	updater          bool
	useDefaultStatic bool
}

func (c *Config) processArguments(args []string) {

	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--api-port":
			if len(args) > i+1 {
				i++
				c.apiPort = args[i]
			}
		case "-p", "--port":
			if len(args) > i+1 {
				i++
				c.sitePort = args[i]
			}
		case "-b", "--base-dir":
			if len(args) > i+1 {
				i++
				c.baseDir = args[i]
			}
		case "--updater":
			c.updater = true
		case "--default-static":
			c.useDefaultStatic = true
		default:
			if !strings.HasPrefix(args[i], "-test.") {
				log.Printf("Ignoring unknown option %s", args[i])
			}
		}

	}
}

func getConfiguration(args []string) *Config {
	c := &Config{}
	c.processArguments(args)

	if c.apiPort == "" {
		c.apiPort = "8081"
	}

	if c.sitePort == "" {
		c.sitePort = "8080"
	}

	if c.baseDir == "" {
		c.baseDir = fmt.Sprintf("%s/.alive", os.Getenv("HOME"))
	}
	c.baseDir = filepath.Clean(c.baseDir)
	c.staticFilePath = filepath.Clean(fmt.Sprintf("%s/static", c.baseDir))
	c.dataFile = filepath.Clean(fmt.Sprintf("%s/data.json", c.baseDir))

	return c
}
