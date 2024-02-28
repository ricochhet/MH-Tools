package manager

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ricochhet/mhwarchivemanager/pkg/config"
	"github.com/ricochhet/mhwarchivemanager/pkg/fsprovider"
	"github.com/ricochhet/mhwarchivemanager/pkg/logger"
	"github.com/ricochhet/mhwarchivemanager/pkg/sevenzip"
)

func T_InstallDirectory(profileName string) error {
	if len(profileName) == 0 {
		profileName = config.DefaultProfileName
	}

	dirPath := fsprovider.Relative(config.DataDirectory, profileName, config.OutputDirectory)
	if err := fsprovider.RemoveAll(dirPath); err != nil {
		logger.SharedLogger.GoRoutineError(err.Error())
		return err
	}

	tempPath := fsprovider.Relative(config.DataDirectory, profileName, config.TempDirectory)
	indexPath := fsprovider.Relative(config.DataDirectory, profileName, config.IndexFile)
	if err := os.MkdirAll(filepath.Dir(indexPath), os.ModePerm); err != nil {
		logger.SharedLogger.GoRoutineError(err.Error())
		return err
	}

	file, err := os.OpenFile(indexPath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		logger.SharedLogger.GoRoutineError(err.Error())
		return err
	}
	defer file.Close()

	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		logger.SharedLogger.GoRoutineError(err.Error())
		return err
	}

	entries, err := sevenzip.ExtractFromList(file, tempPath)
	if err != nil {
		logger.SharedLogger.GoRoutineError(err.Error())
		return err
	}

	if len(entries) != 0 {
		for _, v := range entries {
			if err := filepath.Walk(v, func(walkPath string, info os.FileInfo, err error) error {
				if err != nil {
					logger.SharedLogger.GoRoutineError(err.Error())
					return err
				}
				if info.IsDir() && strings.ToLower(info.Name()) == "nativepc" {
					logger.SharedLogger.Info("Copying nativePC: " + walkPath)
					fsprovider.CopyDirectory(dirPath, walkPath)
				}
				return nil
			}); err != nil {
				logger.SharedLogger.GoRoutineError(err.Error())
				return err
			}
		}
	} else {
		logger.SharedLogger.GoRoutineError("No entries to extract")
		return fmt.Errorf("no entries to extract")
	}

	logger.SharedLogger.Info("Cleaning temp files")
	err = fsprovider.RemoveAll(tempPath)

	if err != nil {
		logger.SharedLogger.GoRoutineError(err.Error())
		return err
	} else {
		logger.SharedLogger.Info("Done")
	}

	return nil
}
