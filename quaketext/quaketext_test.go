package quaketext_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/qw-serverstat/quaketext"
)

func TestToPlainText(t *testing.T) {
	const (
		Letters   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		Numbers   = "0123456789"
		MiscChars = "~!@#$%^&*()-_+={}[]\\|<,>.?/\"';:`"
		Ascii     = Letters + Numbers + MiscChars
	)

	var testedBytes = make(map[byte]bool, 0)

	// white and red ascii
	for _, charByte := range []byte(Ascii) {
		expectedChar := string(charByte)

		// white text
		assert.Equal(t, expectedChar, quaketext.ToPlainText([]byte{charByte}))
		testedBytes[charByte] = true

		// red text
		charByteRedColor := charByte + 128
		assert.Equal(t, expectedChar, quaketext.ToPlainText([]byte{charByteRedColor}))
		testedBytes[charByteRedColor] = true
	}

	// yellow and brown numbers
	for _, charByte := range []byte(Numbers) {
		expectedNumber := string(charByte)

		charByteYellowColor := charByte - 30
		assert.Equal(t, expectedNumber, quaketext.ToPlainText([]byte{charByteYellowColor}))
		testedBytes[charByteYellowColor] = true

		charByteBrownColor := charByteYellowColor + 128
		assert.Equal(t, expectedNumber, quaketext.ToPlainText([]byte{charByteBrownColor}))
		testedBytes[charByteBrownColor] = true
	}

	// colored non-alphanumeric chars
	var miscCharsMap = map[string][]byte{
		" ": {12, 12 + 128, 138},
		"•": {28, 28 + 128, 32, 32 + 128},
		".": {5, 5 + 128, 14, 14 + 128, 15, 15 + 128},
		"<": {29, 127, 128, 157},
		"=": {30, 30 + 128, 129},
		">": {31, 130, 141, 159},
		"[": {16, 16 + 128},
		"]": {17, 17 + 128},
	}

	for expectedChar, specialCharBytes := range miscCharsMap {
		for _, charByte := range specialCharBytes {
			assert.Equal(t, expectedChar, quaketext.ToPlainText([]byte{charByte}), charByte)
			testedBytes[charByte] = true
		}
	}

	// unknown chars
	unknownCharBytes := []byte{0, 1, 2, 3, 4, 6, 7, 8, 9, 10, 11, 13, 131, 132, 134, 135, 136, 137, 139}
	expectedChar := "#"

	for _, charByte := range unknownCharBytes {
		assert.Equal(t, expectedChar, quaketext.ToPlainText([]byte{charByte}), charByte)
		testedBytes[charByte] = true
	}

	// validate test coverage
	for i := byte(0); i < byte(255); i++ {
		if !testedBytes[i] {
			log.Printf("Did not test %d, expected '%s'", i, quaketext.ToPlainText([]byte{i}))
		}
	}

	totalTested := 1 + len(testedBytes)
	expectTested := 255 + 1

	if totalTested < expectTested {
		t.Fatalf("Did not test every character. Tested %d of %d", totalTested, expectTested)
	}
}

func ExampleToPlainText() {
	quakeText := []byte{109, 109, 91, 99, 104, 97, 114, 93, 28, 32, 32, 91, 116, 101, 115, 116, 93, 109, 109}
	plainText := quaketext.ToPlainText(quakeText)
	fmt.Println(plainText)
	// Output: mm[char]•••[test]mm
}
