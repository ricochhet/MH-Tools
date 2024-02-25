@echo off
mkdir build

echo Building MHWArchiveManager GUI
cd cmd/gui
go build -o ../../build/MHWArchiveManager.exe -ldflags="-s -w -H=windowsgui -extldflags=-static" .

echo Building MHWArchiveManager CLI
cd ../cli
go build -o ../../build/MHWArchiveManager.CLI.exe