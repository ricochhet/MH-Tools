package manager

import (
	"fmt"
	"os"

	g "github.com/AllenDang/giu"
	"github.com/ricochhet/mhwarchivemanager/pkg/config"
	"github.com/ricochhet/mhwarchivemanager/pkg/fsprovider"
	"github.com/ricochhet/mhwarchivemanager/pkg/logger"
)

var ui_IndexPath string
var ui_ProfileName string

var ui_ComboItems []string
var ui_ComboSelection int32

func ui_ButtonLaunch() {
	go AGR_LaunchProgram(g.Update)
}

func ui_ButtonExit() {
	go AGR_Exit()
}

func ui_ButtonAdd() {
	if err := A_AddProfile(ui_ProfileName, g.Update); err != nil {
		logger.SharedLogger.Error(err.Error())
		return
	}

	if arr, err := A_UpdateProfileList(); err != nil {
		logger.SharedLogger.Error(err.Error())
		return
	} else {
		ui_ComboItems = arr
	}
}

func ui_ButtonRemove() {
	if err := A_RemoveProfile(ui_ProfileName, g.Update); err != nil {
		logger.SharedLogger.Error(err.Error())
		return
	}

	if arr, err := A_UpdateProfileList(); err != nil {
		logger.SharedLogger.Error(err.Error())
		return
	} else {
		ui_ComboItems = arr
	}
}

func ui_ButtonIndex() {
	go AGR_IndexDirectory(ui_ComboItems[ui_ComboSelection], ui_IndexPath, g.Update)
}

func ui_ButtonInstall() {
	go AGR_InstallDirectory(ui_ComboItems[ui_ComboSelection], g.Update)
}

func ui_MainWindow() {
	g.SingleWindow().Layout(
		g.Label("MHWArchiveManager"),
		g.Row(
			g.Button("Launch Program").OnClick(ui_ButtonLaunch),
			g.Button("Exit MHWArchiveManager").OnClick(ui_ButtonExit),
		),
		g.Row(
			g.Label("Index Path: "),
			g.InputText(&ui_IndexPath).Hint("./path/to/archives"),
			g.Button("Start Index").OnClick(ui_ButtonIndex),
		),
		g.Row(
			g.Label("Profile: "),
			g.InputText(&ui_ProfileName).Hint("Profile Name"),
			g.Button("Add").OnClick(ui_ButtonAdd),
			g.Button("Remove").OnClick(ui_ButtonRemove),
		),
		g.Row(
			g.Label("Profile: "),
			g.Combo("##", ui_ComboItems[ui_ComboSelection], ui_ComboItems, &ui_ComboSelection),
		),
		g.Button("-- Create Installation Folder --").OnClick(ui_ButtonInstall),
		g.Label("Output"),
		g.ListBox("##Logger", logger.LogCache),
	)
}

func A_InitializeUI() {
	if arr, err := A_UpdateProfileList(); err != nil {
		logger.SharedLogger.Error(err.Error())
		return
	} else {
		ui_ComboItems = arr
	}

	wnd := g.NewMasterWindow("MHWArchiveManager", 640, 360, 0)
	wnd.Run(ui_MainWindow)
}

func A_DummyUpdateFunc() {
	logger.SharedLogger.Debug("Called: A_DummyUpdateFunc()")
}

func AGR_LaunchProgram(giuUpdate func()) {
	logger.ClearCache()
	giuUpdate()

	go Launch()
}

func AGR_IndexDirectory(ptr_ComboItem string, ptr_IndexPath string, giuUpdate func()) {
	logger.ClearCache()
	savedIndexPathStr, err := GetSavedIndexPath(ptr_ComboItem)
	if err != nil {
		logger.SharedLogger.GoRoutineError(err.Error())
		return
	}

	if len(ptr_IndexPath) != 0 {
		err := A_HandleIndexDirectory(ptr_ComboItem, ptr_IndexPath, true, giuUpdate)
		if err != nil {
			logger.SharedLogger.GoRoutineError(err.Error())
			return
		}
	}

	if len(ptr_IndexPath) == 0 && len(savedIndexPathStr) != 0 {
		err := A_HandleIndexDirectory(ptr_ComboItem, savedIndexPathStr, false, giuUpdate)
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

func A_HandleIndexDirectory(ptr_ComboItem string, path string, saveIndexPath bool, giuUpdate func()) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
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

func AGR_InstallDirectory(ptr_ComboItem string, giuUpdate func()) {
	logger.ClearCache()
	logger.SharedLogger.Info("Creating installation at path: " + fsprovider.Relative(config.DataDirectory, config.OutputDirectory))
	giuUpdate()

	go InstallDirectory(ptr_ComboItem)
}

func A_AddProfile(ptr_ProfileName string, giuUpdate func()) error {
	logger.ClearCache()
	if len(ptr_ProfileName) == 0 {
		return fmt.Errorf("profile name cannot be blank")
	}

	logger.SharedLogger.Info("Adding profile: " + ptr_ProfileName)
	giuUpdate()

	err := AddProfile(ptr_ProfileName)
	if err != nil {
		return err
	}

	return nil
}

func A_RemoveProfile(ptr_ProfileName string, giuUpdate func()) error {
	logger.ClearCache()
	if len(ptr_ProfileName) == 0 {
		return fmt.Errorf("profile name cannot be blank")
	}

	logger.SharedLogger.Info("Removing profile: " + ptr_ProfileName)
	giuUpdate()

	err := RemoveProfile(ptr_ProfileName)
	if err != nil {
		return err
	}

	return nil
}

func A_UpdateProfileList() ([]string, error) {
	var slice []string
	slice = append(slice, config.DefaultProfileName)

	profileList, err := ReadAllProfiles()
	if err != nil {
		logger.SharedLogger.Error(err.Error())
		return slice, err
	}

	slice = append(slice, profileList...)
	return slice, nil
}

func AGR_Exit() {
	os.Exit(0)
}
