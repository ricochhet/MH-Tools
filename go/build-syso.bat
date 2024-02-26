@echo off

windres MHWArchiveManager.rc -O coff -o MHWArchiveManager.syso
move MHWArchiveManager.syso ./cmd/main

windres MHWArchiveManager.rc -O coff -o MHWArchiveManager.Cli.syso
move MHWArchiveManager.Cli.syso ./cmd/cli_only

windres MHWArchiveManager.rc -O coff -o MHWArchiveManager.Gui.syso
move MHWArchiveManager.Gui.syso ./cmd/gui_only