package sevenzip

import (
	"fmt"

	"github.com/ricochhet/mhwarchivemanager/pkg/logger"
	"github.com/ricochhet/mhwarchivemanager/pkg/process"
)

const (
	NO_ERROR          = 1000
	PROCESS_NOT_FOUND = 1001
	COULD_NOT_EXTRACT = 1002
)

func Extract(source string, destination string) (int, error) {
	if !process.DoesExecutableExist("7z") {
		logger.SharedLogger.Error("7zip was not found.")
		return PROCESS_NOT_FOUND, fmt.Errorf("7zip was not found")
	}

	if err := process.RunExecutable("7z", true, "x", source, "-o"+destination+"/*"); err != nil {
		logger.SharedLogger.Error("Error extracting " + source + ": " + err.Error())
		return COULD_NOT_EXTRACT, err
	}

	return NO_ERROR, nil
}
