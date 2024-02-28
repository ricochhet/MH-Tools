package util

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/ricochhet/mhwarchivemanager/pkg/logger"
)

func WriteQuestGMDLanguages(path string, language string) error {
	if len(path) != 0 {
		languages := []string{"eng", "ara", "chS", "chT", "fre", "ger", "ita", "jpn", "kor", "pol", "ptB", "rus", "spa"}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			return err
		}

		bytes := make([]byte, stat.Size())
		_, err = bufio.NewReader(file).Read(bytes)
		if err != nil && err != io.EOF {
			return err
		}

		for _, lang := range languages {
			newFilePath := strings.Replace(path, language, lang, 1)
			err := os.WriteFile(newFilePath, bytes, 0644)
			if err != nil {
				return err
			}

			logger.SharedLogger.Info("File created: " + newFilePath)
		}
	}

	return nil
}
