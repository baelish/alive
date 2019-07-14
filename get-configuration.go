package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
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

func (c *Config) processArguments() {

	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--api-port":
			if len(os.Args) > i+1 {
				i++
				c.apiPort = os.Args[i]
			}
		case "-p", "--port":
			if len(os.Args) > i+1 {
				i++
				c.sitePort = os.Args[i]
			}
		case "-b", "--base-dir":
			if len(os.Args) > i+1 {
				i++
				c.baseDir = os.Args[i]
			}
		case "--updater":
			c.updater = true
		case "--default-static":
			c.useDefaultStatic = true
		default:
			log.Printf("Ignoring unknown option %s", os.Args[i])
		}

	}
}

func getConfiguration() *Config {
	c := &Config{}
	c.updater = false
	c.processArguments()

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

	var _, err = os.Stat(c.dataFile)
	if os.IsNotExist(err) {
		var file, err = os.Create(c.dataFile)
		if err != nil {
			log.Printf("Data file did not exist and could not create an empty one.")
			log.Fatal(err)
		}
		defer file.Close()
		log.Printf("Created empty data file %s", c.dataFile)
	}
	return c
}
