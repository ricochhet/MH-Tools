package manager

import (
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

func buttonLaunch() {
	go A_LaunchProgram(g.Update)
}

func buttonExit() {
	go A_Exit()
}

func buttonAdd() {
	A_AddProfile(ui_ProfileName, g.Update)
	ui_ComboItems = A_UpdateProfileList()
}

func buttonRemove() {
	A_RemoveProfile(ui_ProfileName, g.Update)
	ui_ComboItems = A_UpdateProfileList()
}

func buttonIndex() {
	go A_IndexDirectory(ui_ComboItems[ui_ComboSelection], ui_IndexPath, g.Update)
}

func buttonInstall() {
	go A_InstallDirectory(ui_ComboItems[ui_ComboSelection], g.Update)
}

func ui() {
	g.SingleWindow().Layout(
		g.Label("MHWArchiveManager"),
		g.Row(
			g.Button("Launch Program").OnClick(buttonLaunch),
			g.Button("Exit MHWArchiveManager").OnClick(buttonExit),
		),
		g.Row(
			g.Label("Index Path: "),
			g.InputText(&ui_IndexPath).Hint("./path/to/archives"),
			g.Button("Start Index").OnClick(buttonIndex),
		),
		g.Row(
			g.Label("Profile: "),
			g.InputText(&ui_ProfileName).Hint("Profile Name"),
			g.Button("Add").OnClick(buttonAdd),
			g.Button("Remove").OnClick(buttonRemove),
		),
		g.Row(
			g.Label("Profile: "),
			g.Combo("##", ui_ComboItems[ui_ComboSelection], ui_ComboItems, &ui_ComboSelection),
		),
		g.Button("-- Create Installation Folder --").OnClick(buttonInstall),
		g.Label("Output"),
		g.ListBox("##Logger", logger.LogCache),
	)
}

func A_InitializeUI() {
	ui_ComboItems = A_UpdateProfileList()
	wnd := g.NewMasterWindow("MHWArchiveManager", 640, 360, 0)
	wnd.Run(ui)
}

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
