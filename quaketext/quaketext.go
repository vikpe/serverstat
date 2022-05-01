package quaketext

/* TODO:
type MarkupFunction func(value string, colorCode rune) string

func ToMarkup(quakeText string, markupFunc MarkupFunction) string {
	return ""
}*/

type QuakeText struct {
	Source []byte
}

func NewFromString(quakeText string) *QuakeText {
	return NewFromBytes([]byte(quakeText))
}

func NewFromBytes(quakeTextBytes []byte) *QuakeText {
	return &QuakeText{Source: quakeTextBytes}
}

func (qtext QuakeText) ToPlainString() string {
	charsetToprows := []string{
		"•", "", "", "", "", "•", "", "", "", "", "", "", "", "", "•", "•",
		"[", "]", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "•", "", "", "",
	}
	charsetTopRowsSize := byte(len(charsetToprows))
	plainTextBytes := make([]byte, 0)

	for _, sourceByte := range qtext.Source {
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
