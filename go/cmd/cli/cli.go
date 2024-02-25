package main

import (
	"os"

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

	if args, err := util.Cmd(os.Args, "compare", 2); err == nil {
		if err := fsprovider.CompareFolders(args[0], args[1]); err != nil {
			logger.SharedLogger.Error(err.Error())
		}
	}

	if args, err := util.Cmd(os.Args, "copy", 2); err == nil {
		if err := fsprovider.CopyDirectory(args[0], args[1]); err != nil {
			logger.SharedLogger.Error(err.Error())
		}
	}

	if args, err := util.Cmd(os.Args, "extract", 2); err == nil {
		if err := process.Extract(args[0], args[1]); err != nil {
			logger.SharedLogger.Error(err.Error())
		}
	}

	if args, err := util.Cmd(os.Args, "delete", 1); err == nil {
		if err := fsprovider.RemoveAll(fsprovider.Relative(args[0])); err != nil {
			logger.SharedLogger.Error(err.Error())
		}
	}

	if args, err := util.Cmd(os.Args, "write_quest_gmd_languages", 0); err == nil {
		if err := util.WriteQuestGMDLanguages(args[0], args[1]); err != nil {
			logger.SharedLogger.Error(err.Error())
		}
	}
}
