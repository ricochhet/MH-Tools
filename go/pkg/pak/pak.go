package pak

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/ricochhet/mhwarchivemanager/pkg/c"
	"github.com/ricochhet/mhwarchivemanager/pkg/fsprovider"
	"github.com/ricochhet/mhwarchivemanager/pkg/logger"
	"github.com/ricochhet/mhwarchivemanager/pkg/thirdparty/murmurhash3"
)

func ProcessDirectory(path string, outputFile string, embed bool) {
	directory, _ := filepath.Abs(path)
	sortedFiles := fsprovider.GetSortedFiles(filepath.Join(directory, "natives"))
	writer, err := c.NewWriter(outputFile, false)

	var data []c.DataEntry
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

		data = append(data, c.DataEntry{Hash: hash, FileName: text})
		data = append(data, c.DataEntry{Hash: hash2, FileName: text})
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

	if embed {
		writer.SeekFromEnd(0)
		WriteData(writer, data)
	} else {
		dataWriter, err := c.NewWriter(outputFile+".data", false)
		if err != nil {
			logger.SharedLogger.Error(err.Error())
			return
		}

		WriteData(dataWriter, data)
		dataWriter.Close()
	}

	writer.Close()
}

func ExtractDirectory(path string, outputDirectory string, embed bool) {
	reader, err := c.NewReader(path)
	var table []c.DataEntry
	if embed {
		table = ReadData(reader)
	} else {
		dataReader, err := c.NewReader(path + ".data")
		if err != nil {
			logger.SharedLogger.Error(err.Error())
			return
		}
		table = ReadData(dataReader)
		dataReader.Close()
	}
	if err != nil {
		logger.SharedLogger.Error(err.Error())
		return
	}

	reader.Seek(0, 0)
	unk0, _ := reader.ReadUInt32()
	unk1, _ := reader.ReadUInt32()
	unk2, _ := reader.ReadUInt32()
	reader.ReadUInt32()

	if unk0 != 1095454795 || unk1 != 4 {
		logger.SharedLogger.Error("Invalid file format")
		return
	}

	var list []c.FileEntry

	for i := uint32(0); i < unk2; i++ {
		fileEntry := c.FileEntry{}
		fileEntry.FileNameLower, _ = reader.ReadUInt32()
		fileEntry.FileNameUpper, _ = reader.ReadUInt32()
		fileEntry.Offset, _ = reader.ReadUInt64()
		fileEntry.UncompSize, _ = reader.ReadUInt64()
		reader.SeekFromCurrent(8)
		reader.SeekFromCurrent(8)
		reader.SeekFromCurrent(4)
		reader.SeekFromCurrent(4)
		list = append(list, fileEntry)
	}

	for _, entry := range list {
		dataEntry := c.FindByHash(table, entry.FileNameLower)
		if dataEntry == nil {
			logger.SharedLogger.Error("File entry not found")
			break
		}

		filePath := filepath.Join(outputDirectory, dataEntry.FileName)
		fileData := make([]byte, entry.UncompSize)
		reader.Read(fileData)

		if len(fileData) > 1073741824 {
			logger.SharedLogger.Error("File too large")
			return
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return
		}

		writer, err := c.NewWriter(filePath, false)
		if err != nil {
			logger.SharedLogger.Error(err.Error())
			return
		}
		writer.Write(fileData)
		writer.Close()
	}

	reader.Close()
}
