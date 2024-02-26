package sevenzip

import (
	"fmt"

	"github.com/ricochhet/mhwarchivemanager/pkg/logger"
	"github.com/ricochhet/mhwarchivemanager/pkg/process"
)

type ErrorCode int

const (
	NoError ErrorCode = iota
	ProcessNotFound
	CouldNotExtract
)

func Extract(source string, destination string) (ErrorCode, error) {
	if !process.DoesExecutableExist("7z") {
		logger.SharedLogger.Error("7zip was not found.")
		return ProcessNotFound, fmt.Errorf("7zip was not found")
	}

	if err := process.RunExecutable("7z", true, "x", source, "-o"+destination+"/*"); err != nil {
		logger.SharedLogger.Error("Error extracting " + source + ": " + err.Error())
		return CouldNotExtract, err
	}

	return NoError, nil
}
