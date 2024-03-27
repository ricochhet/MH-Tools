package manager

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/ricochhet/mhwarchivemanager/pkg/config"
	"github.com/ricochhet/mhwarchivemanager/pkg/fsprovider"
	"github.com/ricochhet/mhwarchivemanager/pkg/logger"
	"github.com/ricochhet/mhwarchivemanager/pkg/sevenzip"
	"github.com/ricochhet/mhwarchivemanager/pkg/thirdparty/copy"
)

func T_InstallDirectory(profileName string) error {
	if len(profileName) == 0 {
		profileName = config.DefaultProfileName
	}

	mtOutputPath := fsprovider.Relative(config.DataDirectory, profileName, config.MtNativePC)
	reOutputPath := fsprovider.Relative(config.DataDirectory, profileName, config.ReNativePC)
	refOutputPath := fsprovider.Relative(config.DataDirectory, profileName, config.RefNativePC)

	if err := fsprovider.RemoveAll(mtOutputPath); err != nil {
		logger.SharedLogger.GoRoutineError(err.Error())
		return err
	}

	if err := fsprovider.RemoveAll(reOutputPath); err != nil {
		logger.SharedLogger.GoRoutineError(err.Error())
		return err
	}

	if err := fsprovider.RemoveAll(refOutputPath); err != nil {
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

	if err := os.MkdirAll(mtOutputPath, os.ModePerm); err != nil {
		logger.SharedLogger.GoRoutineError(err.Error())
		return err
	}

	if err := os.MkdirAll(reOutputPath, os.ModePerm); err != nil {
		logger.SharedLogger.GoRoutineError(err.Error())
		return err
	}

	if err := os.MkdirAll(refOutputPath, os.ModePerm); err != nil {
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

				if info.IsDir() {
					if err := checkFilePaths(info, walkPath, mtOutputPath); err != nil {
						return err
					}

					if err := checkFilePaths(info, walkPath, reOutputPath); err != nil {
						return err
					}

					if err := checkFilePaths(info, walkPath, refOutputPath); err != nil {
						return err
					}
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

func checkFilePaths(info fs.FileInfo, walkPath, outputPath string) error {
	if strings.EqualFold(info.Name(), config.MtNativePC) {
		logger.SharedLogger.Info("Copying nativePC (MT): " + walkPath)
		err := copy.Copy(walkPath, outputPath)
		if err != nil {
			return err
		}
	}

	return nil
}
