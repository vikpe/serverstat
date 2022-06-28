package qstring

import (
	"strings"

	"github.com/vikpe/serverstat/qtext/qchar"
	"github.com/vikpe/serverstat/qutil"
)

type QuakeString string

func New(str string) QuakeString {
	return QuakeString(str)
}

func (qs QuakeString) ToPlainString() string {
	return ToPlainString(string(qs))
}

func (qs QuakeString) ToColorCodes() string {
	return ToColorCodes(string(qs))
}

func (qs QuakeString) MarshalJSON() ([]byte, error) {
	return qutil.MarshalNoEscapeHtml(qs.ToPlainString())
}

func ToPlainString(str string) string {
	plainText := strings.Builder{}

	for _, charByte := range []byte(str) {
		plainText.WriteString(qchar.ToPlainString(charByte))
	}

	return plainText.String()
}

func ToColorCodes(str string) string {
	colorCodes := strings.Builder{}

	for _, charByte := range []byte(str) {
		colorCodes.WriteString(qchar.ToColorCode(charByte))
	}

	return colorCodes.String()
}
