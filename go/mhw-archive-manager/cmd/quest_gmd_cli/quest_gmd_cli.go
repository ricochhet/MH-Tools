package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ricochhet/mhwarchivemanager/pkg/logger"
)

func main() {
	pathPtr := flag.String("path", "", "pathPtr")
	flag.Parse()

	if len(*pathPtr) != 0 {
		languages := []string{"ara", "chS", "chT", "fre", "ger", "ita", "jpn", "kor", "pol", "ptB", "rus", "spa"}
		file, err := os.Open(*pathPtr)
		if err != nil {
			logger.SharedLogger.Error(err.Error())
			return
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			logger.SharedLogger.Error(err.Error())
			return
		}

		bytes := make([]byte, stat.Size())
		_, err = bufio.NewReader(file).Read(bytes)
		if err != nil && err != io.EOF {
			logger.SharedLogger.Error(err.Error())
			return
		}

		for _, lang := range languages {
			newFilePath := strings.Replace(*pathPtr, "eng", lang, 1)
			err := os.WriteFile(newFilePath, bytes, 0644)
			if err != nil {
				logger.SharedLogger.Error("Error writing file: " + err.Error())
				fmt.Printf("Error writing file: %v\n", err)
				return
			}

			logger.SharedLogger.Info("File created: " + newFilePath)
		}
	}
}
