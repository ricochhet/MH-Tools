package pak

import "github.com/ricochhet/mhwarchivemanager/pkg/c"

func WriteData(writer *c.Writer, data []c.DataEntry) {
	startPos, _ := writer.Position()
	for _, entry := range data {
		writer.WriteUInt32(entry.Hash)
		writer.WriteChar(entry.FileName + "\000")
	}
	endPos, _ := writer.Position()
	writer.WriteUInt64(uint64(endPos - startPos))
}

func ReadData(reader *c.Reader) []c.DataEntry {
	var data []c.DataEntry

	reader.SeekFromEnd(-8)
	dataSize, _ := reader.ReadUInt64()
	reader.SeekFromEnd(int64(-dataSize - 8))
	pos, _ := reader.Position()
	size, _ := reader.Size()

	for pos < size-8 {
		pos, _ = reader.Position()
		hash, _ := reader.ReadUInt32()
		var fileName string
		for {
			c, _ := reader.ReadChar()
			if c == '\000' {
				break
			}

			fileName += string(c)
		}

		data = append(data, c.DataEntry{Hash: hash, FileName: fileName})
	}

	return data
}
