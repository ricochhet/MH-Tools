package main

import (
	"os"

	"github.com/ricochhet/mhwarchivemanager/pkg/fsprovider"
	"github.com/ricochhet/mhwarchivemanager/pkg/logger"
	"github.com/ricochhet/mhwarchivemanager/pkg/process"
	"github.com/ricochhet/mhwarchivemanager/pkg/util"
)

func main() {
	arg1 := util.GetStringAtIndex(os.Args, 1)
	arg2 := util.GetStringAtIndex(os.Args, 2)
	arg3 := util.GetStringAtIndex(os.Args, 3)

	if arg1 == "help" {
		logger.SharedLogger.Info("Usage: ")
		logger.SharedLogger.Info(" compare <folder1> <folder2>")
		logger.SharedLogger.Info(" copy <source> <destination>")
		logger.SharedLogger.Info(" extract <source> <destination>")
		logger.SharedLogger.Info(" delete <path>")
	}

	if arg1 == "compare" {
		if err := fsprovider.CompareFolders(arg2, arg3); err != nil {
			logger.SharedLogger.Error(err.Error())
		}
	}

	if arg1 == "copy" {
		if err := fsprovider.CopyDirectory(arg2, arg3); err != nil {
			logger.SharedLogger.Error(err.Error())
		}
	}

	if arg1 == "extract" {
		if err := process.Extract(arg2, arg3); err != nil {
			logger.SharedLogger.Error(err.Error())
		}
	}

	if arg1 == "delete" {
		if err := fsprovider.RemoveAll(fsprovider.Relative(arg2)); err != nil {
			logger.SharedLogger.Error(err.Error())
		}
	}
}
