package qchar

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
