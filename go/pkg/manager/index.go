package manager

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/ricochhet/mhwarchivemanager/pkg/config"
	"github.com/ricochhet/mhwarchivemanager/pkg/fsprovider"
	"github.com/ricochhet/mhwarchivemanager/pkg/logger"
)

func GetSavedIndexPath(profileName string) (string, error) {
	if len(profileName) == 0 {
		profileName = config.DefaultProfileName
	}

	indexPathSaved := fsprovider.Relative(config.DataDirectory, profileName, config.SavedIndexPathFile)
	if err := os.MkdirAll(filepath.Dir(indexPathSaved), os.ModePerm); err != nil {
		return "", err
	}

	file, err := os.OpenFile(indexPathSaved, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer file.Close()

	entries, err := fsprovider.ScanValidEntries(file)
	if err != nil {
		return "", err
	}

	if len(entries) != 0 {
		return entries[0], nil
	}

	return "", fmt.Errorf("no entries in saved index file")
}

func ExcludeFromIndex(profileName string) ([]string, error) {
	if len(profileName) == 0 {
		profileName = config.DefaultProfileName
	}

	exclusionPath := fsprovider.Relative(config.DataDirectory, profileName, config.ExclusionFile)
	if err := os.MkdirAll(filepath.Dir(exclusionPath), os.ModePerm); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(exclusionPath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	entries, err := fsprovider.ScanValidEntries(file)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func T_SaveIndexPath(profileName string, directoryPath string) error {
	if len(profileName) == 0 {
		profileName = config.DefaultProfileName
	}

	indexPathSaved := fsprovider.Relative(config.DataDirectory, profileName, config.SavedIndexPathFile)
	if err := os.MkdirAll(filepath.Dir(indexPathSaved), os.ModePerm); err != nil {
		logger.SharedLogger.GoRoutineError(err.Error())
		return err
	}

	if file, err := os.OpenFile(indexPathSaved, os.O_RDWR|os.O_CREATE, 0666); err == nil {
		defer file.Close()
		fsprovider.Overwrite(file)
		file.WriteString(directoryPath + "\n")
	} else {
		logger.SharedLogger.GoRoutineError(err.Error())
		return err
	}

	return nil
}

func T_IndexDirectory(profileName string, directoryPath string) error {
	if len(profileName) == 0 {
		profileName = config.DefaultProfileName
	}

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

	exclusionEntries, err := ExcludeFromIndex(profileName)
	if err != nil {
		logger.SharedLogger.GoRoutineError(err.Error())
		return err
	}

	existingEntries, err := fsprovider.ScanValidEntries(file)
	if err != nil {
		logger.SharedLogger.GoRoutineError(err.Error())
		return err
	}

	err = filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.SharedLogger.GoRoutineError(err.Error())
			return err
		}

		ext := filepath.Ext(path)
		if slices.Contains(config.ValidFileTypes, ext) && !slices.Contains(existingEntries, path) {
			if slices.Contains(exclusionEntries, path) {
				return nil
			}

			logger.SharedLogger.Info("Adding: " + path)
			existingEntries = append(existingEntries, path)
		}

		return nil
	})

	if err != nil {
		logger.SharedLogger.GoRoutineError(err.Error())
		return err
	}

	fsprovider.Overwrite(file)
	fsprovider.WriteEntriesToFile(file, existingEntries)

	return nil
}
