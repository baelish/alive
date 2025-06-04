package server

import (
	"os"

	"go.uber.org/zap"
)

const emptyDataFile = "[]"

func createStaticContent() {
	if options.Debug {
		logger.Info("Creating Static Content")
	}
	for _, file := range AssetNames() {
		if _, err := os.Stat(options.StaticPath + "/" + file); os.IsNotExist(err) {
			logger.Info("file doesn't exist, creating default file", zap.String("file", options.StaticPath+file))
			err = RestoreAsset(options.StaticPath, file)
			if err != nil {
				logger.Error(err.Error())
			}
		} else if options.DefaultStatic {
			logger.Info("default files enforced, creating default file", zap.String("file", options.StaticPath+file))
			err = RestoreAsset(options.StaticPath, file)
			if err != nil {
				logger.Error(err.Error())
			}
		}
	}
}
