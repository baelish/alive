package main

import "testing"
import "os"
import "fmt"


func TestDefaults(t *testing.T) {
  args := []string{}
  config := getConfiguration(args)
  if config.apiPort != "8081" {
    t.Error ("Expected 8081, got", config.apiPort)
  }
  if config.sitePort != "8080" {
    t.Error ("Expected 8080, got", config.sitePort)
  }
  baseDir := fmt.Sprintf("%s/.alive", os.Getenv("HOME"))
  if config.baseDir != baseDir {
    t.Error ( "Expected ",baseDir, " got ", config.baseDir)
  }
  dataFile := fmt.Sprintf("%s/data.json", baseDir)
  if config.dataFile != dataFile {
    t.Error ( "Expected ",dataFile, " got ", config.dataFile)
  }
  staticFilePath := fmt.Sprintf("%s/static", baseDir)
  if config.staticFilePath != staticFilePath {
    t.Error ( "Expected ",staticFilePath, " got ", config.staticFilePath)
  }
  if config.updater != false {
    t.Error ( "Expected false got ", config.updater )
  }
  if config.useDefaultStatic != false {
    t.Error ( "Expected false got ", config.useDefaultStatic )
  }
}

func TestArgumentProcessing(t *testing.T) {
  config := &Config{}
  args := []string{"","-b", "/data", "-p", "1234"}
  config.processArguments(args)
  if config.baseDir != "/data" {
    t.Error ("Expected /data got",config.baseDir)
  }
  if config.sitePort != "1234" {
    t.Error ("Expected 1234 got",config.sitePort)
  }
  args = []string{"","--api-port", "1233", "--port", "1235", "--base-dir", "/var/data", "--updater", "--default-static" }
  config.processArguments(args)
  if config.apiPort != "1233" {
    t.Error ("Expected 1233 got",config.apiPort)
  }
  if config.baseDir != "/var/data" {
    t.Error ("Expected /var/data got",config.baseDir)
  }
  if config.sitePort != "1235" {
    t.Error ("Expected 1235 got",config.sitePort)
  }
  if config.updater != true {
    t.Error ( "Expected true got ", config.updater )
  }
  if config.useDefaultStatic != true {
    t.Error ( "Expected true got ", config.useDefaultStatic )
  }
}
