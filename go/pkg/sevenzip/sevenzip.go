package sevenzip

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/ricochhet/mhwarchivemanager/pkg/fsprovider"
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
		return ProcessNotFound, fmt.Errorf("7zip was not found")
	}

	if err := process.RunExecutable("7z", true, "x", source, "-o"+destination+"/*"); err != nil {
		return CouldNotExtract, err
	}

	return NoError, nil
}

func ExtractFromList(file *os.File, tempPath string) ([]string, error) {
	extractedDirs := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if len(scanner.Text()) == 0 {
			continue
		}
		zipFilePath := strings.TrimSpace(scanner.Text())
		extractedDir := path.Join(tempPath, fsprovider.FileNameWithoutExtension(zipFilePath))
		logger.SharedLogger.Info("Extracting: " + extractedDir)
		szerr, err := Extract(zipFilePath, tempPath)
		if szerr == ProcessNotFound {
			break
		}

		if szerr == CouldNotExtract {
			logger.SharedLogger.Error("Error extracting " + zipFilePath + ": " + err.Error())
			continue
		}

		extractedDirs = append(extractedDirs, extractedDir)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return extractedDirs, nil
}
