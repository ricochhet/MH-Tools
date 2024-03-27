package sevenzip

import (
	"fmt"

	"github.com/ricochhet/mhwarchivemanager/pkg/fsprovider"
	"github.com/ricochhet/mhwarchivemanager/pkg/logger"
	"github.com/ricochhet/mhwarchivemanager/pkg/thirdparty/copy"
)

func builtinExtractFn(args []string) {
	if err := checkArgs(args, 2); err != nil {
		return
	}

	if _, err := Extract(args[0], args[1]); err != nil {
		logger.SharedLogger.Error(err.Error())
	}
}

func builtinCopyFn(args []string) {
	if err := checkArgs(args, 2); err != nil {
		return
	}

	if err := copy.Copy(args[0], args[1]); err != nil {
		logger.SharedLogger.Error(err.Error())
	}
}

func builtinDeleteFn(args []string) {
	if err := checkArgs(args, 1); err != nil {
		return
	}

	if err := fsprovider.RemoveAll(fsprovider.Relative(args[0])); err != nil {
		logger.SharedLogger.Error(err.Error())
	}
}

func checkArgs(args []string, expected int) error {
	if len(args) != expected {
		s := fmt.Sprintf("expected: %d arguments but got %d", expected, len(args))
		logger.SharedLogger.Error(s)
		return fmt.Errorf(s)
	}
	return nil
}
