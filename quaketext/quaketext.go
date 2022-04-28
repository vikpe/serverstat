package quaketext

/* TODO:
type MarkupFunction func(value string, colorCode rune) string

func ToMarkup(quakeText string, markupFunc MarkupFunction) string {
	return ""
}*/

func BytesToPlainString(quakeTextBytes []byte) string {
	charsetToprows := []string{
		"•", "", "", "", "", "•", "", "", "", "", "", "", "", "", "•", "•",
		"[", "]", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "•", "", "", "",
	}
	charsetTopRowsSize := byte(len(charsetToprows))
	plainTextBytes := make([]byte, 0)

	for _, sourceByte := range quakeTextBytes {
		sourceByte &= 0x7f // remove color

		if sourceByte == 127 { // weird left arrow at end of charset
			continue
		} else if sourceByte < charsetTopRowsSize {
			translatedBytes := []byte(charsetToprows[sourceByte])
			plainTextBytes = append(plainTextBytes, translatedBytes...)
		} else {
			plainTextBytes = append(plainTextBytes, sourceByte)
		}
	}

	return string(plainTextBytes)
}

func StringToPlainString(quakeText string) string {
	return BytesToPlainString([]byte(quakeText))
}
