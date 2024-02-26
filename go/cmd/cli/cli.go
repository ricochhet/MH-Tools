package main

import (
	"os"
	"strconv"

	"github.com/ricochhet/mhwarchivemanager/pkg/fsprovider"
	"github.com/ricochhet/mhwarchivemanager/pkg/logger"
	"github.com/ricochhet/mhwarchivemanager/pkg/manager"
	"github.com/ricochhet/mhwarchivemanager/pkg/pak"
	"github.com/ricochhet/mhwarchivemanager/pkg/sevenzip"
	"github.com/ricochhet/mhwarchivemanager/pkg/util"
)

var CommandList = map[string]string{
	"Usage:":                                        "",
	"  compare <folder1> <folder2>":                 "",
	"  copy <source> <destination>":                 "",
	"  extract <source> <destination>":              "",
	"  delete <path>":                               "",
	"Pak Tools:":                                    "",
	"  pak <folder> <output> <0/1[embed data]>":     "",
	"  unpak <folder> <output> <0/1[embed data]>":   "",
	"  compress <file>":                             "",
	"  decompress <file>":                           "",
	"Misc Tools:":                                   "",
	"  write_quest_gmd_languages <path> <language>": "",
}

func main() {
	if _, err := util.Cmd(os.Args, "help", 0); err == nil {
		for command, description := range CommandList {
			logger.SharedLogger.Info(command)
			if description != "" {
				logger.SharedLogger.Info("  " + description)
			}
		}
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
		manager.A_LaunchProgram(manager.A_DummyUpdateFunc)
	}

	if args, err := util.Cmd(os.Args, "index", 2); err == nil {
		manager.A_IndexDirectory(args[0], args[1], manager.A_DummyUpdateFunc)
	}

	if args, err := util.Cmd(os.Args, "install", 2); err == nil {
		manager.A_InstallDirectory(args[0], manager.A_DummyUpdateFunc)
	}

	if args, err := util.Cmd(os.Args, "add_profile", 1); err == nil {
		manager.A_AddProfile(args[0], manager.A_DummyUpdateFunc)
	}

	if args, err := util.Cmd(os.Args, "remove_profile", 1); err == nil {
		manager.A_RemoveProfile(args[0], manager.A_DummyUpdateFunc)
	}

	if args, err := util.Cmd(os.Args, "write_quest_gmd_languages", 0); err == nil {
		if err := util.WriteQuestGMDLanguages(args[0], args[1]); err != nil {
			logger.SharedLogger.Error(err.Error())
		}
	}

	if args, err := util.Cmd(os.Args, "pak", 3); err == nil {
		selection, err := strconv.Atoi(args[2])
		if err != nil {
			return
		}
		pak.ProcessDirectory(args[0], args[1], selection != 0)
	}

	if args, err := util.Cmd(os.Args, "unpak", 3); err == nil {
		selection, err := strconv.Atoi(args[2])
		if err != nil {
			return
		}
		pak.ExtractDirectory(args[0], args[1], selection != 0)
	}

	if args, err := util.Cmd(os.Args, "compress", 1); err == nil {
		pak.CompressPakData(args[0])
	}

	if args, err := util.Cmd(os.Args, "decompress", 1); err == nil {
		pak.DecompressPakData(args[0])
	}
}
