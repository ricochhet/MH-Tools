package sevenzip

import (
	"fmt"

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
		return ProcessNotFound, fmt.Errorf("7zip was not found")
	}

	if err := process.RunExecutable("7z", true, "x", source, "-o"+destination+"/*"); err != nil {
		return CouldNotExtract, err
	}

	return NoError, nil
}
