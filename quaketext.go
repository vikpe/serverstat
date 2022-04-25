package serverstat

import "strings"

func quakeTextToPlainText(quakeText string) string {
	readableTextBytes := []byte(quakeText)

	var charset = [...]byte{
		' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
		'[', ']', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', ' ', ' ', ' ', ' ',
	}

	for i := range quakeText {
		readableTextBytes[i] &= 0x7f

		if quakeText[i] < byte(len(charset)) {
			readableTextBytes[i] = charset[quakeText[i]]
		}
	}

	return strings.TrimSpace(string(readableTextBytes))
}
