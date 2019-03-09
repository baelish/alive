package main

import (
    "path/filepath"
    "fmt"
    "log"
    "os"
)

type Config struct {
    baseDir string
    staticFilePath string
}

func (c *Config) processArguments() {

    for i := 1; i < len(os.Args); i++ {
        switch os.Args[i] {
        case "-b", "--base-dir":
          if len(os.Args) > i+1 {
            i++
            c.baseDir = os.Args[i]
          }
        default:
            log.Printf("Ignoring unknown option %s", os.Args[i])
        }

    }
}


func getConfiguration() (*Config){
    c := &Config{}
    c.processArguments()
    if c.baseDir == "" {
        c.baseDir = fmt.Sprintf("%s/.alive",os.Getenv("HOME"))
    }
    c.baseDir = filepath.Clean(c.baseDir)
    c.staticFilePath = filepath.Clean(fmt.Sprintf("%s/static",c.baseDir))
    return c
}
