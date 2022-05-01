package quaketext_test

import (
	"fmt"

	"github.com/vikpe/qw-serverstat/quaketext"
)

func ExampleNewFromString() {
	quakeText := "XantoM"
	qtext := quaketext.NewFromString(quakeText)
	fmt.Println(qtext.ToPlainString())
	// Output: XantoM
}

func ExampleNewFromBytes() {
	quakeTextBytes := []byte{88, 97, 110, 116, 111, 77}
	qtext := quaketext.NewFromBytes(quakeTextBytes)
	fmt.Println(qtext.ToPlainString())
	// Output: XantoM
}

func ExampleQuakeText_ToPlainString() {
	quakeText := "XantoM"
	plainString := quaketext.NewFromString(quakeText).ToPlainString()
	fmt.Println(plainString)
}
