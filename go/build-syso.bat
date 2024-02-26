@echo off

windres MHWArchiveManager.rc -O coff -o MHWArchiveManager.syso
move MHWArchiveManager.syso ./cmd/main

windres MHWArchiveManager.rc -O coff -o MHWArchiveManager.CLI.syso
move MHWArchiveManager.CLI.syso ./cmd/cli_only