package quaketext_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/qw-serverstat/quaketext"
)

func TestToPlainText(t *testing.T) {
	// printable range of ascii table
	var ascii []byte

	for i := ' '; i <= '~'; i++ {
		ascii = append(ascii, byte(i))
	}

	// keep track of tested bytes
	testedBytes := make(map[byte]bool, 128)

	for b := 0; b < 128; b++ {
		testedBytes[byte(b)] = false
	}

	// test ascii table
	for i := range ascii {
		charByte := ascii[i]
		char := string([]byte{charByte})

		// normal/white ascii
		charWhite := char
		assert.Equal(t, char, quaketext.StringToPlainString(charWhite))
		testedBytes[charByte] = true

		// red ascii
		charByteRed := charByte + 128
		charRed := string([]byte{charByteRed})
		assert.Equal(t, char, quaketext.StringToPlainString(charRed))
		testedBytes[charByte+128] = true

		// yellow numbers
		if char >= "0" && char <= "9" {
			charByteYellow := charByte - 30
			charYellow := string([]byte{charByteYellow})
			assert.Equal(t, char, quaketext.StringToPlainString(charYellow)) // yellow numbers
			testedBytes[charByteYellow] = true
		}
	}

	// test top two rows of charset + last char (127)
	specialChars := map[string][]byte{
		"":  {1, 2, 3, 4, 6, 7, 8, 9, 10, 11, 12, 13, 29, 30, 31, 127},
		"•": {0, 5, 14, 15, 28},
		"[": {16},
		"]": {17},
	}

	for expectedChar, charBytes := range specialChars {
		for _, charByte := range charBytes {
			chars := string([]byte{charByte})
			assert.Equal(t, expectedChar, quaketext.StringToPlainString(chars), charByte)
			testedBytes[charByte] = true
		}
	}

	// validate that all bytes are tested
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

func ExampleStringToPlainText() {
	quakeText := "XantoM"
	plainText := quaketext.StringToPlainString(quakeText)
	fmt.Println(plainText)
	// Output: XantoM
}

func ExampleBytesToPlainText() {
	quakeTextBytes := []byte{88, 97, 110, 116, 111, 77}
	plainText := quaketext.BytesToPlainString(quakeTextBytes)
	fmt.Println(plainText)
	// Output: XantoM
}
