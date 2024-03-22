package pak

import (
	"github.com/ricochhet/mhwarchivemanager/pkg/c"
	"github.com/ricochhet/mhwarchivemanager/pkg/fsprovider"
	"github.com/ricochhet/mhwarchivemanager/pkg/logger"
)

func CompressPakData(path string) {
	writer, err := c.NewWriter(path, true)
	if err != nil {
		logger.SharedLogger.Error(err.Error())
		return
	}

	reader, err := c.NewReader(path + ".data")
	if err != nil {
		logger.SharedLogger.Error(err.Error())
		return
	}

	data := ReadData(reader)
	writer.SeekFromEnd(0)
	WriteData(writer, data)

	reader.Close()
	writer.Close()

	if err := fsprovider.RemoveAll(fsprovider.Relative(path + ".data")); err != nil {
		logger.SharedLogger.Error(err.Error())
		return
	}
}

func DecompressPakData(path string) {
	reader, err := c.NewReader(path)
	if err != nil {
		logger.SharedLogger.Error(err.Error())
		return
	}

	writer, err := c.NewWriter(path+".data", false)
	if err != nil {
		logger.SharedLogger.Error(err.Error())
		return
	}

	data := ReadData(reader)
	WriteData(writer, data)
	writer.Close()

	reader.SeekFromEnd(-8)
	dataSize, _ := reader.ReadUInt64()

	reader.Seek(0, 0)
	size, _ := reader.Size()
	decompSize := size - (int64(dataSize) - 8)
	buffer := make([]byte, decompSize)
	reader.Read(buffer)
	reader.Close()

	decomp, err := c.NewWriter(path, false)
	if err != nil {
		logger.SharedLogger.Error(err.Error())
		return
	}
	decomp.Write(buffer)
	decomp.Close()
}
