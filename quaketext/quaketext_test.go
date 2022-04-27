package quaketext_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/qw-serverstat/quaketext"
)

func TestToPlainText(t *testing.T) {
	// printable ascii table
	var ascii []byte

	for i := ' '; i <= '~'; i++ {
		ascii = append(ascii, byte(i))
	}

	testedBytes := make(map[byte]bool, 128)

	for b := 0; b < 128; b++ {
		testedBytes[byte(b)] = false
	}

	for i := range ascii {
		charByte := ascii[i]
		char := string(charByte)

		// normal ascii
		charsWhite := string([]byte{charByte})
		assert.Equal(t, char, quaketext.ToPlainText(charsWhite))

		// red ascii
		charByteRed := charByte + 128
		charsRed := string([]byte{charByteRed})
		assert.Equal(t, char, quaketext.ToPlainText(charsRed))

		// yellow numbers
		if char >= "0" && char <= "9" {
			charsYellow := string([]byte{charByte - 30})
			assert.Equal(t, char, quaketext.ToPlainText(charsYellow)) // yellow numbers

			testedBytes[charByte-30] = true
		}

		testedBytes[charByte] = true
		testedBytes[charByte+128] = true
	}

	// top two rows of charset + last char (127)
	specialCases := map[string][]byte{
		"":  {1, 2, 3, 4, 6, 7, 8, 9, 10, 11, 12, 13, 29, 30, 31, 127},
		"•": {0, 5, 14, 15, 28},
		"[": {16},
		"]": {17},
	}

	for expectedChar, charBytes := range specialCases {
		for _, charByte := range charBytes {
			chars := string([]byte{charByte})
			assert.Equal(t, expectedChar, quaketext.ToPlainText(chars), charByte)

			testedBytes[charByte] = true
		}
	}

	hasTestedAllBytes := true

	for byte_, value := range testedBytes {
		if !value {
			log.Println("did not test", byte_)
			hasTestedAllBytes = false
		}
	}

	if !hasTestedAllBytes {
		t.Fatal("Did not test all chars.")
	}
}

func ExampleToPlainText() {
	quakeText := "XantoM"
	plainText := quaketext.ToPlainText(quakeText)
	fmt.Println(plainText)
	// Output: XantoM
}
