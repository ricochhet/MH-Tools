package sevenzip

import (
	"fmt"
	"os"
	"strings"

	"github.com/ricochhet/mhwarchivemanager/pkg/fsprovider"
	"github.com/ricochhet/mhwarchivemanager/pkg/logger"
	"github.com/ricochhet/mhwarchivemanager/pkg/thirdparty/copy"
)

func Parse(path string) error {
	if !fsprovider.Exists(path) {
		return fmt.Errorf("file: " + path + " does not exist")
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	entries, err := fsprovider.Scanner(file)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if strings.HasPrefix(entry, "#") {
			continue
		}

		parts := splitArgs(entry)
		if len(parts) < 2 {
			return fmt.Errorf("not enough arguments provided to function")
		}

		fnName := strings.ToLower(parts[0])
		fnArgs := parts[1:]
		switch fnName {
		case "extract":
			builtinExtractFn(fnArgs)
		case "copy":
			builtinCopyFn(fnArgs)
		case "delete":
			builtinDeleteFn(fnArgs)
		default:
		}
	}

	return nil
}

func splitArgs(input string) []string {
	var parts []string
	var part strings.Builder
	inQuote := false

	for _, char := range input {
		if char == '"' {
			inQuote = !inQuote
		} else if char == ' ' && !inQuote {
			parts = append(parts, part.String())
			part.Reset()
		} else {
			part.WriteRune(char)
		}
	}

	parts = append(parts, part.String())
	return parts
}

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
