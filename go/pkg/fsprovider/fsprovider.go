package fsprovider

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ricochhet/mhwarchivemanager/pkg/logger"
)

func WriteEntriesToFile(file *os.File, entries []string) {
	for _, entry := range entries {
		file.WriteString(entry + "\n")
	}
}

func ScanValidEntries(file *os.File) ([]string, error) {
	existingEntries := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if len(scanner.Text()) == 0 {
			continue
		}
		if _, err := os.Stat(scanner.Text()); err == nil {
			existingEntries = append(existingEntries, scanner.Text())
		} else {
			logger.SharedLogger.Info("Invalid Entry: " + scanner.Text())
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return existingEntries, nil
}

func Relative(directories ...string) string {
	result := "./" + directories[0]

	for _, dir := range directories[1:] {
		result = path.Join(result, dir)
	}

	return result
}

func Overwrite(file *os.File) error {
	if err := file.Truncate(0); err != nil {
		return err
	}

	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	return nil
}

func FileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName))
}

func RemoveAll(fileName string) error {
	err := os.RemoveAll(fileName)
	if err != nil {
		return err
	}

	return nil
}

func CopyDirectory(destination, source string) error {
	return filepath.Walk(source, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		outpath := filepath.Join(destination, strings.TrimPrefix(path, source))

		if info.IsDir() {
			os.MkdirAll(outpath, info.Mode())
			return nil
		}

		if !info.Mode().IsRegular() {
			switch info.Mode().Type() & os.ModeType {
			case os.ModeSymlink:
				link, err := os.Readlink(path)
				if err != nil {
					return err
				}
				return os.Symlink(link, outpath)
			}
			return nil
		}

		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()

		fh, err := os.Create(outpath)
		if err != nil {
			return err
		}
		defer fh.Close()

		fh.Chmod(info.Mode())

		_, err = io.Copy(fh, in)
		return err
	})
}

func CalculateChecksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func CompareFolders(folder1, folder2 string) error {
	fileChecksums1 := make(map[string]string)
	fileChecksums2 := make(map[string]string)

	err := filepath.Walk(folder1, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relativePath, _ := filepath.Rel(folder1, path)
			checksum, err := CalculateChecksum(path)
			if err != nil {
				return err
			}
			fileChecksums1[relativePath] = checksum
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = filepath.Walk(folder2, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relativePath, _ := filepath.Rel(folder2, path)
			checksum, err := CalculateChecksum(path)
			if err != nil {
				return err
			}
			fileChecksums2[relativePath] = checksum
		}
		return nil
	})
	if err != nil {
		return err
	}

	for path1, checksum1 := range fileChecksums1 {
		checksum2, exists := fileChecksums2[path1]
		if !exists {
			logger.SharedLogger.Info("File " + path1 + " exists in folder1 but not in folder2")
			continue
		}

		if checksum1 != checksum2 {
			logger.SharedLogger.Info("Checksums for file " + path1 + " do not match:")
			logger.SharedLogger.Info("	Folder1: " + checksum1)
			logger.SharedLogger.Info("	Folder2: " + checksum2)
		}
	}

	for path2 := range fileChecksums2 {
		_, exists := fileChecksums1[path2]
		if !exists {
			logger.SharedLogger.Info("File " + path2 + " exists in folder2 but not in folder1")
		}
	}

	return nil
}

func SortByParentAndName(files []string) []string {
	sort.Slice(files, func(i, j int) bool {
		parentA := filepath.Dir(files[i])
		parentB := filepath.Dir(files[j])

		if parentA == parentB {
			return filepath.Base(files[i]) < filepath.Base(files[j])
		}
		return parentA < parentB
	})

	return files
}

func GetSortedFiles(directory string) []string {
	var files []string
	filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	return SortByParentAndName(files)
}
