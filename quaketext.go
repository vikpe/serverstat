package serverstat

import (
	"strings"
)

func quakeTextToPlainText(quakeText string) string {
	var charset = [128]string{
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

	translation := ""

	for _, charByte := range []byte(quakeText) {
		charByte &= 0x7f // strip color (> 128)
		translation += charset[charByte]
	}

	return strings.TrimSpace(translation)
}
