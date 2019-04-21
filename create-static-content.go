package main

import (
	"log"
	"os"
)

func createStaticContent(path string) {
	for _, file := range AssetNames() {
		if _, err := os.Stat(path + "/" + file); os.IsNotExist(err) {
			log.Printf("'%s/%s' doesn't exist, creating default file.", path, file)
			RestoreAsset(path, file)
		} else if config.useDefaultStatic {
			log.Printf("Default files enforced, creating default file '%s/%s'.", path, file)
			RestoreAsset(path, file)
		}
	}
}
