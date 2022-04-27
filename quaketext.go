package serverstat

import (
	"strings"
)

func quakeTextToPlainText(quakeText string) string {
	const charsetSize = 128
	var charset = [charsetSize]string{
		"<", "=", ">", "#", "#", ".", "#", "#",
		"#", "#", " ", "#", " ", ">", ".", ".",
		"[", "]", "0", "1", "2", "3", "4", "5",
		"6", "7", "8", "9", "•", "<", "=", ">",
		"•", "!", "\"", "#", "$", "%", "&", "\"",
		"(", ")", "*", "+", ",", "-", ".", "/",
		"0", "1", "2", "3", "4", "5", "6", "7",
		"8", "9", ":", ";", "<", "=", ">", "?",
		"@", "A", "B", "C", "D", "E", "F", "G",
		"H", "I", "J", "K", "L", "M", "N", "O",
		"P", "Q", "R", "S", "T", "U", "V", "W",
		"X", "Y", "Z", "[", "\\", "]", "^", "_",
		"`", "a", "b", "c", "d", "e", "f", "g",
		"h", "i", "j", "k", "l", "m", "n", "o",
		"p", "q", "r", "s", "t", "u", "v", "w",
		"x", "y", "z", "{", "|", "}", "~", "<",
	}

	plainText := ""

	for _, charByte := range []byte(quakeText) {
		charByte &= 0x7f // strip color (> 128)

		if charByte < charsetSize {
			plainText += charset[charByte]
		} else {
			plainText += "?"
		}
	}

	return strings.TrimSpace(plainText)
}
