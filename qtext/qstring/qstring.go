package qstring

import (
	"github.com/vikpe/serverstat/qtext/qchar"
)

func ToPlainString(str string) string {
	plainText := ""

	for _, charByte := range []byte(str) {
		plainText += qchar.ToPlainString(charByte)
	}

	return plainText
}
