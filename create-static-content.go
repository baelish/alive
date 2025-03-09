package main

import (
	"log"
	"os"
)

const emptyDataFile = "[]"

func createStaticContent() {
	if options.Debug {
		log.Printf("Creating Static Content")
	}
	for _, file := range AssetNames() {
		if _, err := os.Stat(options.StaticPath + "/" + file); os.IsNotExist(err) {
			log.Printf("'%s/%s' doesn't exist, creating default file.", options.StaticPath, file)
			err = RestoreAsset(options.StaticPath, file)
			if err != nil {
				log.Print(err)
			}
		} else if options.DefaultStatic {
			log.Printf("Default files enforced, creating default file '%s/%s'.", options.StaticPath, file)
			err = RestoreAsset(options.StaticPath, file)
			if err != nil {
				log.Print(err)
			}
		}
	}
}
