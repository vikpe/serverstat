package serverstat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestHelloName(t *testing.T) {
	//foo := "he he"
	//nameStr := []byte(foo)

	// a-z, A-Z
	/*for char := '0'; char < 'z'; char++ {
		lowerAndUpper := fmt.Sprintf("%c%c", char, unicode.ToUpper(char))
		log.Println(lowerAndUpper, int(char), int(unicode.ToUpper(char)))
		assert.Equal(t, lowerAndUpper, quakeTextToPlainText(lowerAndUpper))
	}*/

	name := []byte{109, 109, 91, 99, 104, 97, 114, 93, 28, 32, 32, 91, 116, 101, 115, 116, 93, 109, 109}

	assert.Equal(t, "mm[char]•••[test]mm", quakeTextToPlainText(string(name)))
}
