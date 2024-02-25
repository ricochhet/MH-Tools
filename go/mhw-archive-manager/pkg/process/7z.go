package process

import (
	"fmt"

	"github.com/ricochhet/mhwarchivemanager/pkg/logger"
)

func Extract(source string, destination string) error {
	if !DoesExecutableExist("7z") {
		logger.SharedLogger.Error("7zip was not found.")
		return fmt.Errorf("7zip was not found")
	}

	if err := RunExecutable("7z", true, "x", source, "-o"+destination+"/*"); err != nil {
		logger.SharedLogger.Error("Error extracting " + source + ": " + err.Error())
		return err
	}

	return nil
}
