package main

import (
	"flag"
	"os"

	"github.com/ricochhet/mhwarchivemanager/pkg/config"
	"github.com/ricochhet/mhwarchivemanager/pkg/fsprovider"
	"github.com/ricochhet/mhwarchivemanager/pkg/logger"
	"github.com/ricochhet/mhwarchivemanager/pkg/manager"
)

func main() {
	indexPathPtr := flag.String("index", "", "indexPathPtr")
	installPathPtr := flag.Bool("install", false, "installPathPtr")
	launchPtr := flag.Bool("launch", false, "launchPtr")
	flag.Parse()

	if len(*indexPathPtr) != 0 {
		logger.SharedLogger.Info("Creating index")
		err := manager.IndexDirectory(config.DefaultProfileName, *indexPathPtr)
		if err != nil {
			logger.SharedLogger.Error(err.Error())
		}
	}

	if *installPathPtr {
		logger.SharedLogger.Info("Creating install")
		err := manager.InstallDirectory(config.DefaultProfileName)
		if err != nil {
			logger.SharedLogger.Error(err.Error())
		}
	}

	if *launchPtr {
		logger.SharedLogger.Info("Launching program")
		err := manager.Launch()
		if err != nil {
			logger.SharedLogger.Error(err.Error())
		}
	}

	if os.Args[1] == "diff" {
		folder1 := os.Args[2]
		folder2 := os.Args[3]

		err := fsprovider.CompareFolders(folder1, folder2)
		if err != nil {
			logger.SharedLogger.Error(err.Error())
		}
	}
}
