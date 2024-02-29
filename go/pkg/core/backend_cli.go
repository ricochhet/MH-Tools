package core

import (
	"os"
	"strconv"

	"github.com/ricochhet/mhwarchivemanager/pkg/fsprovider"
	"github.com/ricochhet/mhwarchivemanager/pkg/logger"
	"github.com/ricochhet/mhwarchivemanager/pkg/pak"
	"github.com/ricochhet/mhwarchivemanager/pkg/sevenzip"
	"github.com/ricochhet/mhwarchivemanager/pkg/util"
)

func A_InitializeCommandLine() {
	if _, err := util.Cmd(os.Args, "help", 0); err == nil {
		logger.SharedLogger.Info("Usage: ")
		logger.SharedLogger.Info("  compare <folder1> <folder2>")
		logger.SharedLogger.Info("  copy <source> <destination>")
		logger.SharedLogger.Info("  extract <source> <destination>")
		logger.SharedLogger.Info("  delete <path>")
		logger.SharedLogger.Info("Pak Tools: ")
		logger.SharedLogger.Info("  pak <folder> <output> <0/1[embed data]>")
		logger.SharedLogger.Info("  unpak <folder> <output> <0/1[embed data]>")
		logger.SharedLogger.Info("  compress <file>")
		logger.SharedLogger.Info("  decompress <file>")
		logger.SharedLogger.Info("Misc Tools: ")
		logger.SharedLogger.Info("  write_quest_gmd_languages <path> <language>")
	}

	if args, err := util.Cmd(os.Args, "compare", 2); err == nil {
		if err := fsprovider.CompareFolders(args[0], args[1]); err != nil {
			logger.SharedLogger.Error(err.Error())
		}
	}

	if args, err := util.Cmd(os.Args, "checksum", 1); err == nil {
		checksum, err := fsprovider.CalculateChecksum(args[0])
		if err != nil {
			logger.SharedLogger.Error(err.Error())
			return
		}

		logger.SharedLogger.Info(checksum)
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
		go A_LaunchProgramThread(A_DummyUpdateFunc)
	}

	if args, err := util.Cmd(os.Args, "index", 2); err == nil {
		go A_IndexDirectoryThread(args[0], args[1], A_DummyUpdateFunc)
	}

	if args, err := util.Cmd(os.Args, "install", 2); err == nil {
		go A_InstallDirectoryThread(args[0], A_DummyUpdateFunc)
	}

	if args, err := util.Cmd(os.Args, "add_profile", 1); err == nil {
		A_AddProfile(args[0], A_DummyUpdateFunc)
	}

	if args, err := util.Cmd(os.Args, "remove_profile", 1); err == nil {
		A_RemoveProfile(args[0], A_DummyUpdateFunc)
	}

	if args, err := util.Cmd(os.Args, "write_quest_gmd_languages", 0); err == nil {
		if err := util.WriteQuestGMDLanguages(args[0], args[1]); err != nil {
			logger.SharedLogger.Error(err.Error())
		}
	}

	if args, err := util.Cmd(os.Args, "pak", 3); err == nil {
		selection, err := strconv.Atoi(args[2])
		if err != nil {
			logger.SharedLogger.Error(err.Error())
			return
		}
		pak.ProcessDirectory(args[0], args[1], selection != 0)
	}

	if args, err := util.Cmd(os.Args, "unpak", 3); err == nil {
		selection, err := strconv.Atoi(args[2])
		if err != nil {
			logger.SharedLogger.Error(err.Error())
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
