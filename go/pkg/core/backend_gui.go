package core

import (
	"fmt"

	g "github.com/AllenDang/giu"
	"github.com/ricochhet/mhwarchivemanager/pkg/config"
	"github.com/ricochhet/mhwarchivemanager/pkg/logger"
	"github.com/ricochhet/mhwarchivemanager/pkg/manager"
)

var ui_IndexPath string
var ui_ProfileName string

var ui_ComboItems []string
var ui_ComboSelection int32

func ui_ButtonLaunch() {
	go A_LaunchProgramThread(g.Update)
}

func ui_ButtonExit() {
	go A_ExitThread()
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
	go A_IndexDirectoryThread(ui_ComboItems[ui_ComboSelection], ui_IndexPath, g.Update)
}

func ui_ButtonInstall() {
	go A_InstallDirectoryThread(ui_ComboItems[ui_ComboSelection], g.Update)
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

func A_AddProfile(ptr_ProfileName string, giuUpdate func()) error {
	logger.ClearCache()
	if len(ptr_ProfileName) == 0 {
		return fmt.Errorf("profile name cannot be blank")
	}

	logger.SharedLogger.Info("Adding profile: " + ptr_ProfileName)
	giuUpdate()

	err := manager.AddProfile(ptr_ProfileName)
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

	err := manager.RemoveProfile(ptr_ProfileName)
	if err != nil {
		return err
	}

	return nil
}

func A_UpdateProfileList() ([]string, error) {
	var slice []string
	slice = append(slice, config.DefaultProfileName)

	profileList, err := manager.ReadAllProfiles()
	if err != nil {
		logger.SharedLogger.Error(err.Error())
		return slice, err
	}

	slice = append(slice, profileList...)
	return slice, nil
}
