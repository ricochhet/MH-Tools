package pak

import (
	"encoding/binary"
	"fmt"
	"math"
	"path/filepath"
	"strings"

	"github.com/ricochhet/mhwarchivemanager/pkg/c"
	"github.com/ricochhet/mhwarchivemanager/pkg/fsprovider"
	"github.com/ricochhet/mhwarchivemanager/pkg/logger"
	"github.com/ricochhet/mhwarchivemanager/pkg/murmurhash3"
)

func ProcessDirectory(path string, outputFile string) {
	directory, _ := filepath.Abs(path)
	sortedFiles := fsprovider.GetSortedFiles(filepath.Join(directory, "natives"))
	writer, err := c.NewWriter(outputFile, false)
	var list []c.FileEntry
	if err != nil {
		logger.SharedLogger.Error(err.Error())
		return
	}

	writer.WriteUInt32(1095454795)
	writer.WriteUInt32(4)
	writer.WriteUInt32(uint32(len(sortedFiles)))
	writer.WriteUInt32(0)
	pos, _ := writer.Position()
	writer.Seek(int64(48*len(sortedFiles))+pos, 0)

	for _, obj := range sortedFiles {
		fileEntry2 := c.FileEntry{}
		text := strings.ReplaceAll(obj, directory, "")
		text = strings.ReplaceAll(text, "\\", "/")

		if text[0] == '/' {
			text = text[1:]
		}

		hashBytes := murmurhash3.NewX86_32(math.MaxUint32)
		hashBytes.Write(c.Utf8ToUtf16(strings.ToLower(text)))
		hash := binary.LittleEndian.Uint32(hashBytes.Sum(nil))

		hashBytes2 := murmurhash3.NewX86_32(math.MaxUint32)
		hashBytes2.Write(c.Utf8ToUtf16(strings.ToUpper(text)))
		hash2 := binary.LittleEndian.Uint32(hashBytes2.Sum(nil))

		reader, err := c.NewReader(obj)
		if err != nil {
			logger.SharedLogger.Error(err.Error())
			return
		}

		size, _ := reader.Size()
		array2 := make([]byte, size)
		reader.Read(array2)

		fileEntry2.FileName = text
		pos, _ = writer.Position()
		fileEntry2.Offset = uint64(pos)
		fileEntry2.UncompSize = uint64(len(array2))
		fileEntry2.FileNameLower = uint32(hash)
		fileEntry2.FileNameUpper = uint32(hash2)
		list = append(list, fileEntry2)
		writer.Write(array2)
	}

	writer.SeekFromBeginning(16)
	for _, item := range list {
		fmt.Printf("%s, %v, %v\n", item.FileName, item.FileNameLower, item.FileNameUpper)
		writer.WriteUInt32(item.FileNameLower)
		writer.WriteUInt32(item.FileNameUpper)
		writer.WriteUInt64(item.Offset)
		writer.WriteUInt64(item.UncompSize)
		writer.WriteUInt64(item.UncompSize)
		writer.WriteUInt64(0)
		writer.WriteUInt32(0)
		writer.WriteUInt32(0)
	}

	writer.Close()
}
