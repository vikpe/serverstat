package qchar

import "bytes"

const (
	ColorWhite = "w"
	ColorBrown = "b"
	ColorGold  = "g"
)

func RemoveColor(qchar byte) byte {
	return qchar & 0x7f
}

func ToPlainString(qchar byte) string {
	plainByte := RemoveColor(qchar)

	if plainByte == 127 { // weird left arrow at end of charset
		return ""
	}

	charsetToprows := []string{
		"•", "", "", "", "", "•", "", "", "", "", "", "", "", "", "•", "•",
		"[", "]", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "•", "", "", "",
	}
	charsetTopRowsSize := byte(len(charsetToprows))

	if plainByte < charsetTopRowsSize {
		return charsetToprows[plainByte]
	} else {
		return string(plainByte)
	}
}

func ToColorCode(qchar byte) string {
	rowInCharset := qchar / 16

	if rowInCharset > 9 {
		return ColorBrown
	}

	goldBytes := []byte{
		16, 16 + 128, 17 + 128, // braces
		5 + 128, 14 + 128, 15 + 128, 28 + 128, // dots
	}

	if hasChar(goldBytes, qchar) {
		return ColorGold
	}

	plainChar := RemoveColor(qchar)

	if plainChar >= 18 && plainChar <= 27 { // brown numbers
		return ColorBrown
	}

	return ColorWhite
}

func hasChar(haystack []byte, needle byte) bool {
	return bytes.IndexByte(haystack, needle) != -1
}
