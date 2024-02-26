package main

import (
	"os"

	"github.com/ricochhet/mhwarchivemanager/pkg/fsprovider"
	"github.com/ricochhet/mhwarchivemanager/pkg/logger"
	"github.com/ricochhet/mhwarchivemanager/pkg/manager"
	"github.com/ricochhet/mhwarchivemanager/pkg/pak"
	"github.com/ricochhet/mhwarchivemanager/pkg/sevenzip"
	"github.com/ricochhet/mhwarchivemanager/pkg/util"
)

func _DummyUpdateFunc() {
	logger.SharedLogger.Debug("Called: _DummyUpdateFunc()")
}

func main() {
	if _, err := util.Cmd(os.Args, "help", 0); err == nil {
		logger.SharedLogger.Info("Usage: ")
		logger.SharedLogger.Info(" compare <folder1> <folder2>")
		logger.SharedLogger.Info(" copy <source> <destination>")
		logger.SharedLogger.Info(" extract <source> <destination>")
		logger.SharedLogger.Info(" delete <path>")
		logger.SharedLogger.Info("Tools: ")
		logger.SharedLogger.Info(" write_quest_gmd_languages <path> <language>")
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
		if _, err := sevenzip.Extract(args[0], args[1]); err != nil {
			logger.SharedLogger.Error(err.Error())
		}
	}

	if args, err := util.Cmd(os.Args, "delete", 1); err == nil {
		if err := fsprovider.RemoveAll(fsprovider.Relative(args[0])); err != nil {
			logger.SharedLogger.Error(err.Error())
		}
	}

	if _, err := util.Cmd(os.Args, "launch", 0); err == nil {
		manager.A_LaunchProgram(_DummyUpdateFunc)
	}

	if args, err := util.Cmd(os.Args, "index", 2); err == nil {
		manager.A_IndexDirectory(args[0], args[1], _DummyUpdateFunc)
	}

	if args, err := util.Cmd(os.Args, "install", 2); err == nil {
		manager.A_InstallDirectory(args[0], _DummyUpdateFunc)
	}

	if args, err := util.Cmd(os.Args, "add_profile", 1); err == nil {
		manager.A_AddProfile(args[0], _DummyUpdateFunc)
	}

	if args, err := util.Cmd(os.Args, "remove_profile", 1); err == nil {
		manager.A_RemoveProfile(args[0], _DummyUpdateFunc)
	}

	if args, err := util.Cmd(os.Args, "write_quest_gmd_languages", 0); err == nil {
		if err := util.WriteQuestGMDLanguages(args[0], args[1]); err != nil {
			logger.SharedLogger.Error(err.Error())
		}
	}

	if args, err := util.Cmd(os.Args, "pak", 1); err == nil {
		pak.ProcessDirectory(args[0], "re_chunk_00x.pak.patch_00y.pak", true)
	}

	if args, err := util.Cmd(os.Args, "unpak", 2); err == nil {
		pak.ExtractDirectory(args[0], args[1], true)
	}
}
