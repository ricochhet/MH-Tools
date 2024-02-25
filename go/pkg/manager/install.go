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
	tempPath := fsprovider.Relative(config.DataDirectory, config.SettingsDirectory, profileName, config.TempDirectory)

	if err := fsprovider.RemoveAll(dirPath); err != nil {
		return err
	}

	file, err := os.Open(fsprovider.Relative(config.DataDirectory, config.SettingsDirectory, profileName, config.IndexFile))
	if err != nil {
		logger.SharedLogger.Error(err.Error())
		return err
	}
	defer file.Close()

	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		logger.SharedLogger.Error(err.Error())
		return err
	}

	if entries := extract(file, tempPath); len(entries) != 0 {
		for _, v := range entries {
			if err := filepath.Walk(v, func(walkPath string, info os.FileInfo, err error) error {
				if err != nil {
					logger.SharedLogger.Error(err.Error())
					return err
				}
				if info.IsDir() && strings.ToLower(info.Name()) == "nativepc" {
					logger.SharedLogger.Info("Copying nativePC: " + walkPath)
					fsprovider.CopyDirectory(dirPath, walkPath)
				}
				return nil
			}); err != nil {
				logger.SharedLogger.Error(err.Error())
				return err
			}
		}
	}

	logger.SharedLogger.Info("Cleaning temp files")
	err = fsprovider.RemoveAll(tempPath)

	if err != nil {
		logger.SharedLogger.Error(err.Error())
		return err
	} else {
		logger.SharedLogger.Info("Done")
	}

	return nil
}

func extract(file *os.File, tempPath string) []string {
	extractedDirs := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if len(scanner.Text()) == 0 {
			continue
		}
		zipFilePath := strings.TrimSpace(scanner.Text())
		extractedDir := path.Join(tempPath, fsprovider.FileNameWithoutExtension(zipFilePath))
		logger.SharedLogger.Info("Extracting: " + extractedDir)
		sz, err := sevenzip.Extract(zipFilePath, tempPath)
		if sz == sevenzip.PROCESS_NOT_FOUND {
			logger.SharedLogger.Error("Error " + err.Error())
			break
		}

		if sz == sevenzip.COULD_NOT_EXTRACT {
			logger.SharedLogger.Error("Error extracting " + zipFilePath + ": " + err.Error())
			continue
		}

		extractedDirs = append(extractedDirs, extractedDir)
	}

	if err := scanner.Err(); err != nil {
		logger.SharedLogger.Error(err.Error())
		return nil
	}

	return extractedDirs
}
