@echo off
mkdir build

echo Building MHWArchiveManager
cd cmd/main
@REM go build -o ../../build/MHWArchiveManager.exe -ldflags="-s -w -H=windowsgui -extldflags=-static" .
go build -o ../../build/MHWArchiveManager.exe -ldflags="-s -w -extldflags=-static" .

echo Building MHWArchiveManager.CLI
cd ../cli_only
go build -o ../../build/MHWArchiveManager.CLI.exe -ldflags="-s -w -extldflags=-static" .