package c

import (
	"unicode/utf16"
	"unicode/utf8"
)

func Utf8ToUtf16(utf8str string) string {
	utf8Bytes := []byte(utf8str)
	utf16Runes := make([]uint16, 0, len(utf8Bytes))

	for len(utf8Bytes) > 0 {
		r, size := utf8.DecodeRune(utf8Bytes)
		if r == utf8.RuneError && size == 1 {
			return ""
		}
		utf16Chars := utf16.Encode([]rune{r})
		utf16Runes = append(utf16Runes, utf16Chars...)
		utf8Bytes = utf8Bytes[size:]
	}

	return string(utf16.Decode(utf16Runes))
}
