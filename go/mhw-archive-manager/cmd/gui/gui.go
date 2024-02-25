package main

import (
	g "github.com/AllenDang/giu"
	"github.com/ricochhet/mhwarchivemanager/pkg/logger"
	"github.com/ricochhet/mhwarchivemanager/pkg/manager"
)

var ui_IndexPath string
var ui_ProfileName string
var consoleOutput string

var ui_ComboItems []string
var ui_ComboSelection int32

func buttonLaunch() {
	go manager.A_LaunchProgram(g.Update)
}

func buttonExit() {
	go manager.A_Exit()
}

func buttonAdd() {
	manager.A_AddProfile(ui_ProfileName, g.Update)
	ui_ComboItems = manager.A_UpdateProfileList()
}

func buttonRemove() {
	manager.A_RemoveProfile(ui_ProfileName, g.Update)
	ui_ComboItems = manager.A_UpdateProfileList()
}

func buttonIndex() {
	go manager.A_IndexDirectory(ui_ComboItems[ui_ComboSelection], ui_IndexPath, g.Update)
}

func buttonInstall() {
	go manager.A_InstallDirectory(ui_ComboItems[ui_ComboSelection], g.Update)
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

func main() {
	ui_ComboItems = manager.A_UpdateProfileList()
	wnd := g.NewMasterWindow("MHWArchiveManager", 640, 360, 0)
	wnd.Run(ui)
}
