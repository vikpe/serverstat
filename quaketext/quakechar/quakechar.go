package quakechar

type QuakeChar struct {
	Byte byte
}

func New(b byte) *QuakeChar {
	return &QuakeChar{Byte: b}
}

func (qchar QuakeChar) RemoveColor() *QuakeChar {
	return New(qchar.Byte & 0x7f)
}

func (qchar QuakeChar) ToPlainString() string {
	charsetToprows := []string{
		"•", "", "", "", "", "•", "", "", "", "", "", "", "", "", "•", "•",
		"[", "]", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "•", "", "", "",
	}
	charsetTopRowsSize := byte(len(charsetToprows))

	plainByte := qchar.RemoveColor().Byte

	if plainByte == 127 { // weird left arrow at end of charset
		return ""
	} else if plainByte < charsetTopRowsSize {
		return charsetToprows[plainByte]
	} else {
		return string(plainByte)
	}
}
