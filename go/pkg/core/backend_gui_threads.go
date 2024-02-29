package core

import (
	"os"

	"github.com/ricochhet/mhwarchivemanager/pkg/config"
	"github.com/ricochhet/mhwarchivemanager/pkg/fsprovider"
	"github.com/ricochhet/mhwarchivemanager/pkg/logger"
	"github.com/ricochhet/mhwarchivemanager/pkg/manager"
)

func A_ExitThread() {
	os.Exit(0)
}

func A_LaunchProgramThread(giuUpdate func()) {
	logger.ClearCache()
	giuUpdate()

	go manager.T_Launch()
}

func A_InstallDirectoryThread(ptr_ComboItem string, giuUpdate func()) {
	logger.ClearCache()
	logger.SharedLogger.Info("Creating installation at path: " + fsprovider.Relative(config.DataDirectory, config.OutputDirectory))
	giuUpdate()

	go manager.T_InstallDirectory(ptr_ComboItem)
}

func A_IndexDirectoryThread(ptr_ComboItem string, ptr_IndexPath string, giuUpdate func()) {
	logger.ClearCache()
	savedIndexPathStr, err := manager.GetSavedIndexPath(ptr_ComboItem)
	if err != nil {
		logger.SharedLogger.GoRoutineError(err.Error())
	}

	if len(ptr_IndexPath) != 0 {
		err := handleIndexDirectory(ptr_ComboItem, ptr_IndexPath, true, giuUpdate)
		if err != nil {
			logger.SharedLogger.GoRoutineError(err.Error())
			return
		}
	}

	if len(ptr_IndexPath) == 0 && len(savedIndexPathStr) != 0 {
		err := handleIndexDirectory(ptr_ComboItem, savedIndexPathStr, false, giuUpdate)
		if err != nil {
			logger.SharedLogger.GoRoutineError(err.Error())
			return
		}
	}

	if len(ptr_IndexPath) == 0 && len(savedIndexPathStr) == 0 {
		logger.SharedLogger.GoRoutineError("No index path specified")
		return
	}
}

func handleIndexDirectory(ptr_ComboItem string, path string, saveIndexPath bool, giuUpdate func()) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}

	logger.SharedLogger.Info("Indexing path: " + path)
	giuUpdate()
	go manager.T_IndexDirectory(ptr_ComboItem, path)
	if saveIndexPath {
		go manager.T_SaveIndexPath(ptr_ComboItem, path)
	}

	return nil
}
