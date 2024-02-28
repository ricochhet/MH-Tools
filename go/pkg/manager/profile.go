package manager

import (
	"bufio"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/ricochhet/mhwarchivemanager/pkg/config"
	"github.com/ricochhet/mhwarchivemanager/pkg/fsprovider"
	"github.com/ricochhet/mhwarchivemanager/pkg/logger"
)

func AddProfile(name string) error {
	listFilePath := fsprovider.Relative(config.DataDirectory, config.ProfileListFile)
	if err := os.MkdirAll(filepath.Dir(listFilePath), os.ModePerm); err != nil {
		return err
	}

	file, err := os.OpenFile(listFilePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	items := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if len(scanner.Text()) == 0 {
			continue
		}
		items = append(items, strings.TrimSpace(scanner.Text()))
	}

	if !slices.Contains(items, name) && len(name) != 0 {
		items = append(items, strings.TrimSpace(name))
	}

	fsprovider.Overwrite(file)
	fsprovider.WriteEntriesToFile(file, items)

	if err := scanner.Err(); err != nil {
		return err
	}

	return err
}

func RemoveProfile(name string) error {
	listFilePath := fsprovider.Relative(config.DataDirectory, config.ProfileListFile)
	if err := os.MkdirAll(filepath.Dir(listFilePath), os.ModePerm); err != nil {
		return err
	}

	file, err := os.OpenFile(listFilePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	items := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		trimmed := strings.TrimSpace(scanner.Text())
		if len(scanner.Text()) == 0 || name == trimmed {
			logger.SharedLogger.Info("Removing: " + scanner.Text())
			continue
		}
		items = append(items, trimmed)
	}

	fsprovider.Overwrite(file)
	fsprovider.WriteEntriesToFile(file, items)

	if err := scanner.Err(); err != nil {
		return err
	}

	return err
}

func ReadAllProfiles() ([]string, error) {
	listFilePath := fsprovider.Relative(config.DataDirectory, config.ProfileListFile)
	if err := os.MkdirAll(filepath.Dir(listFilePath), os.ModePerm); err != nil {
		return nil, err
	}
	file, err := os.OpenFile(listFilePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	items := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if len(scanner.Text()) == 0 {
			continue
		}
		items = append(items, strings.TrimSpace(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return items, err
}
