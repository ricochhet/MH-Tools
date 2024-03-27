package sevenzip

import (
	"fmt"
	"os"
	"strings"

	"github.com/ricochhet/mhwarchivemanager/pkg/fsprovider"
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
		Execute(fnName, fnArgs)
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
