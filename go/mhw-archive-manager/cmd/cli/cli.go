package main

import (
	"os"
	"strings"

	"github.com/ricochhet/mhwarchivemanager/pkg/fsprovider"
	"github.com/ricochhet/mhwarchivemanager/pkg/logger"
	"github.com/ricochhet/mhwarchivemanager/pkg/process"
	"github.com/ricochhet/mhwarchivemanager/pkg/util"
)

func main() {
	if _, err := util.Cmd(os.Args, "help", 0); err == nil {
		logger.SharedLogger.Info("Usage: ")
		logger.SharedLogger.Info(" compare <folder1> <folder2>")
		logger.SharedLogger.Info(" copy <source> <destination>")
		logger.SharedLogger.Info(" extract <source> <destination>")
		logger.SharedLogger.Info(" delete <path>")
	}

	if arr, err := util.Cmd(os.Args, "testcmd", 2); err == nil {
		logger.SharedLogger.Info(strings.Join(arr, ", "))
	}

	if arr, err := util.Cmd(os.Args, "compare", 2); err == nil {
		if err := fsprovider.CompareFolders(arr[0], arr[1]); err != nil {
			logger.SharedLogger.Error(err.Error())
		}
	}

	if arr, err := util.Cmd(os.Args, "copy", 2); err == nil {
		if err := fsprovider.CopyDirectory(arr[0], arr[1]); err != nil {
			logger.SharedLogger.Error(err.Error())
		}
	}

	if arr, err := util.Cmd(os.Args, "extract", 2); err == nil {
		if err := process.Extract(arr[0], arr[1]); err != nil {
			logger.SharedLogger.Error(err.Error())
		}
	}

	if arr, err := util.Cmd(os.Args, "delete", 1); err == nil {
		if err := fsprovider.RemoveAll(fsprovider.Relative(arr[0])); err != nil {
			logger.SharedLogger.Error(err.Error())
		}
	}
}
