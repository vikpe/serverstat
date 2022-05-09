package qchar_test

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qtext/qchar"
)

func TestToPlainString(t *testing.T) {
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
		assert.Equal(t, char, qchar.ToPlainString(charByte))
		testedBytes[charByte] = true

		// red ascii
		charByteRed := charByte + 128
		assert.Equal(t, char, qchar.ToPlainString(charByteRed))
		testedBytes[charByteRed] = true

		// yellow numbers
		if char >= "0" && char <= "9" {
			charByteYellow := charByte - 30
			assert.Equal(t, char, qchar.ToPlainString(charByteYellow)) // yellow numbers
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
			assert.Equal(t, expectedChar, qchar.ToPlainString(charByte), charByte)
			testedBytes[charByte] = true
		}
	}

	// validate that all bytes are tested
	hasTestedAllBytes := true

	for byte_, value := range testedBytes {
		if !value {
			log.Println("did not test charbyte:", byte_)
			hasTestedAllBytes = false
		}
	}

	if !hasTestedAllBytes {
		t.Fatal("Did not test all chars.")
	}
}

func TestToColorCode(t *testing.T) {
	// brown
	brownRowStart := 10
	brownRowStop := 15

	for b := brownRowStart * 16; b <= brownRowStop*16; b++ {
		assert.Equal(t, "b", qchar.ToColorCode(byte(b)), b)
	}

	// gold
	goldBytes := []byte{
		16, 16 + 128, 17, 17 + 128, // braces
		5 + 128, 14 + 128, 15 + 128, 28 + 128, // dots
	}

	for _, b := range goldBytes {
		assert.Equal(t, "g", qchar.ToColorCode(b), b)
	}

	// green numbers
	for b := byte(18); b <= 27; b++ {
		assert.Equal(t, "b", qchar.ToColorCode(b), b)
	}
	for b := byte(18) + 128; b <= 27+128; b++ {
		assert.Equal(t, "b", qchar.ToColorCode(b), b)
	}

	// white
	whiteRowStart := 2
	whiteRowStop := 7

	for b := whiteRowStart * 16; b <= whiteRowStop*16; b++ {
		assert.Equal(t, "w", qchar.ToColorCode(byte(b)), b)
	}

	// white dots
	whiteDots := []byte{5, 14, 15, 28}

	for _, b := range whiteDots {
		assert.Equal(t, "w", qchar.ToColorCode(b), b)
	}

	// misc non-printable chars
	assert.Equal(t, "w", qchar.ToColorCode(1), 1)
}

func TestRemoveColor(t *testing.T) {
	assert.Equal(t, byte(100), qchar.RemoveColor(228))
	assert.Equal(t, byte(100), qchar.RemoveColor(100))
}
