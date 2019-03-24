package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// Config struct contains configuration to be used throughout the program
type Config struct {
	baseDir        string
	staticFilePath string
	updater				 bool
}

func (c *Config) processArguments() {

	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-b", "--base-dir":
			if len(os.Args) > i+1 {
				i++
				c.baseDir = os.Args[i]
			}
		case "--updater":
				c.updater = true
		default:
			log.Printf("Ignoring unknown option %s", os.Args[i])
		}

	}
}

func getConfiguration() *Config {
	c := &Config{}
	c.updater = false
	c.processArguments()
	if c.baseDir == "" {
		c.baseDir = fmt.Sprintf("%s/.alive", os.Getenv("HOME"))
	}
	c.baseDir = filepath.Clean(c.baseDir)
	c.staticFilePath = filepath.Clean(fmt.Sprintf("%s/static", c.baseDir))
	return c
}
