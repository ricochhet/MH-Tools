package manager

import (
	"bufio"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/ricochhet/mhwarchivemanager/pkg/config"
	"github.com/ricochhet/mhwarchivemanager/pkg/fsprovider"
	"github.com/ricochhet/mhwarchivemanager/pkg/logger"
	"github.com/ricochhet/mhwarchivemanager/pkg/sevenzip"
)

func InstallDirectory(profileName string) error {
	if len(profileName) == 0 {
		profileName = config.DefaultProfileName
	}

	dirPath := fsprovider.Relative(config.DataDirectory, config.SettingsDirectory, profileName, config.OutputDirectory)
	if err := fsprovider.RemoveAll(dirPath); err != nil {
		logger.SharedLogger.GoRoutineError(err.Error())
		return err
	}

	tempPath := fsprovider.Relative(config.DataDirectory, config.SettingsDirectory, profileName, config.TempDirectory)
	indexPath := fsprovider.Relative(config.DataDirectory, config.SettingsDirectory, profileName, config.IndexFile)
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

	entries, err := extract(file, tempPath)
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
		logger.SharedLogger.Info("No entries to extract")
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

func extract(file *os.File, tempPath string) ([]string, error) {
	extractedDirs := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if len(scanner.Text()) == 0 {
			continue
		}
		zipFilePath := strings.TrimSpace(scanner.Text())
		extractedDir := path.Join(tempPath, fsprovider.FileNameWithoutExtension(zipFilePath))
		logger.SharedLogger.Info("Extracting: " + extractedDir)
		szerr, err := sevenzip.Extract(zipFilePath, tempPath)
		if err != nil {
			return nil, err
		}

		if szerr == sevenzip.ProcessNotFound {
			return nil, err
		}

		if szerr == sevenzip.CouldNotExtract {
			logger.SharedLogger.Error("Error extracting " + zipFilePath + ": " + err.Error())
			continue
		}

		extractedDirs = append(extractedDirs, extractedDir)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return extractedDirs, nil
}
