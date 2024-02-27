@echo off
mkdir build

echo Building MHWArchiveManager
cd cmd/main
go build -o ../../build/MHWArchiveManager.exe -ldflags="-s -w -H=windowsgui -extldflags=-static" .

echo Building MHWArchiveManager.Gui
cd ../gui_only
go build -o ../../build/MHWArchiveManager.Gui.exe -ldflags="-s -w -H=windowsgui -extldflags=-static" .

echo Building MHWArchiveManager.Cli
cd ../cli_only
go build -o ../../build/MHWArchiveManager.Cli.exe -ldflags="-s -w -extldflags=-static" .