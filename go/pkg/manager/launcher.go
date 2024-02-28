package manager

import (
	"errors"
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
		logger.SharedLogger.GoRoutineError(err.Error())
		return err
	}

	file, err := os.OpenFile(launchPath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		logger.SharedLogger.GoRoutineError(err.Error())
		return err
	}
	defer file.Close()

	entries, err := fsprovider.ScanValidEntries(file)
	if err != nil {
		logger.SharedLogger.GoRoutineError(err.Error())
		return err
	}

	if len(entries) != 0 {
		firstEntry := entries[0]
		if _, err := os.Stat(firstEntry); errors.Is(err, os.ErrNotExist) {
			logger.SharedLogger.GoRoutineError(err.Error())
			return err
		}

		err = process.RunExecutable(firstEntry, false)
		if err != nil {
			logger.SharedLogger.GoRoutineError(err.Error())
			return err
		}
	}

	return nil
}
