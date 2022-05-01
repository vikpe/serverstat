package quaketext

import "github.com/vikpe/qw-serverstat/quaketext/quakechar"

/* TODO:
type MarkupFunction func(value string, colorCode rune) string

func ToMarkup(quakeText string, markupFunc MarkupFunction) string {
	return ""
}*/

type QuakeText struct {
	Chars []quakechar.QuakeChar
}

func NewFromString(quakeText string) *QuakeText {
	return NewFromBytes([]byte(quakeText))
}

func NewFromBytes(quakeTextBytes []byte) *QuakeText {
	chars := make([]quakechar.QuakeChar, 0)

	for _, n := range quakeTextBytes {
		chars = append(chars, *quakechar.New(n))
	}

	return &QuakeText{Chars: chars}
}

func (qtext QuakeText) ToPlainString() string {
	plainText := ""

	for _, char := range qtext.Chars {
		plainText += char.ToPlainString()
	}

	return plainText
}
