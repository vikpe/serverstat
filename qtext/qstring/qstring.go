package qstring

import (
	"github.com/vikpe/serverstat/qtext/qchar"
)

type QuakeString struct {
	bytes []byte
}

func New(str string) QuakeString {
	return QuakeString{bytes: []byte(str)}
}

func (qs QuakeString) ToPlainString() string {
	return ToPlainString(string(qs.bytes))
}

func (qs QuakeString) ToColorCodes() string {
	return ToColorCodes(string(qs.bytes))
}

func ToPlainString(str string) string {
	plainText := ""

	for _, charByte := range []byte(str) {
		plainText += qchar.ToPlainString(charByte)
	}

	return plainText
}

func ToColorCodes(str string) string {
	colorCodes := ""

	for _, charByte := range []byte(str) {
		code := qchar.ToColorCode(charByte)
		colorCodes += code
	}

	return colorCodes
}
