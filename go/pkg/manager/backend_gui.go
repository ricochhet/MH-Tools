package manager

import (
	"os"

	"github.com/ricochhet/mhwarchivemanager/pkg/config"
	"github.com/ricochhet/mhwarchivemanager/pkg/fsprovider"
	"github.com/ricochhet/mhwarchivemanager/pkg/logger"
)

func A_DummyUpdateFunc() {
	logger.SharedLogger.Debug("Called: A_DummyUpdateFunc()")
}

func A_LaunchProgram(giuUpdate func()) {
	logger.ClearCache()
	giuUpdate()

	go Launch()
}

func A_IndexDirectory(ptr_ComboItem string, ptr_IndexPath string, giuUpdate func()) {
	logger.ClearCache()
	savedIndexPathStr := GetSavedIndexPath(ptr_ComboItem)

	if len(ptr_IndexPath) != 0 {
		err := A_HandleIndexDirectory(ptr_ComboItem, ptr_IndexPath, true, giuUpdate)
		if err != nil {
			return
		}
	}

	if len(ptr_IndexPath) == 0 && len(savedIndexPathStr) != 0 {
		err := A_HandleIndexDirectory(ptr_ComboItem, savedIndexPathStr, false, giuUpdate)
		if err != nil {
			return
		}
	}

	if len(ptr_IndexPath) == 0 && len(savedIndexPathStr) == 0 {
		logger.SharedLogger.Error("No index path specified")
		return
	}
}

func A_HandleIndexDirectory(ptr_ComboItem string, path string, saveIndexPath bool, giuUpdate func()) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		logger.SharedLogger.Error("Index path does not exist")
		return err
	}

	logger.SharedLogger.Info("Indexing path: " + path)
	giuUpdate()
	go IndexDirectory(ptr_ComboItem, path)
	if saveIndexPath {
		go SaveIndexPath(ptr_ComboItem, path)
	}

	return nil
}

func A_InstallDirectory(ptr_ComboItem string, giuUpdate func()) {
	logger.ClearCache()
	logger.SharedLogger.Info("Creating installation at path: " + fsprovider.Relative(config.DataDirectory, config.OutputDirectory))
	giuUpdate()

	go InstallDirectory(ptr_ComboItem)
}

func A_AddProfile(ptr_ProfileName string, giuUpdate func()) {
	logger.ClearCache()
	if len(ptr_ProfileName) == 0 {
		logger.SharedLogger.Warn("Profile name cannot be blank")
		return
	}

	logger.SharedLogger.Info("Adding profile: " + ptr_ProfileName)
	giuUpdate()

	AddProfile(ptr_ProfileName)
}

func A_RemoveProfile(ptr_ProfileName string, giuUpdate func()) {
	logger.ClearCache()
	if len(ptr_ProfileName) == 0 {
		logger.SharedLogger.Warn("Profile name cannot be blank")
		return
	}

	logger.SharedLogger.Info("Removing profile: " + ptr_ProfileName)
	giuUpdate()

	RemoveProfile(ptr_ProfileName)
}

func A_UpdateProfileList() []string {
	var slice []string
	slice = append(slice, config.DefaultProfileName)

	profileList, err := ReadAllProfiles()
	if err != nil {
		logger.SharedLogger.Error(err.Error())
		return slice
	}

	slice = append(slice, profileList...)
	return slice
}

func A_Exit() {
	os.Exit(0)
}
