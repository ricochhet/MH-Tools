package manager

import (
	"os"
	"path/filepath"

	"github.com/ricochhet/mhwarchivemanager/pkg/config"
	"github.com/ricochhet/mhwarchivemanager/pkg/fsprovider"
	"github.com/ricochhet/mhwarchivemanager/pkg/logger"
	"github.com/ricochhet/mhwarchivemanager/pkg/process"
)

func Launch() error {
	launchPath := fsprovider.Relative(config.DataDirectory, config.SettingsDirectory, config.LaunchFile)
	if err := os.MkdirAll(filepath.Dir(launchPath), os.ModePerm); err != nil {
		return err
	}

	file, err := os.OpenFile(launchPath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	if entries := fsprovider.ScanValidEntries(file); len(entries) != 0 {
		firstEntry := entries[0]
		if !process.DoesExecutableExist(firstEntry) {
			logger.SharedLogger.Warn("The executable " + firstEntry + " could not be found")
			return err
		}

		err = process.RunExecutable(firstEntry, false)
		if err != nil {
			logger.SharedLogger.Error(err.Error())
			return err
		}
	}

	return err
}
