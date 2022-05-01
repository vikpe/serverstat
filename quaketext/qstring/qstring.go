package qstring

import "github.com/vikpe/qw-serverstat/quaketext/qchar"

func ToPlainString(str string) string {
	plainText := ""

	for _, char := range str {
		plainText += qchar.ToPlainString(byte(char))
	}

	return plainText
}
