package qstring_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/qw-serverstat/quaketext/qstring"
)

func TestToPlainString(t *testing.T) {
	expect := "XantuM"
	actual := qstring.ToPlainString("XantõM")
	assert.Equal(t, expect, actual)
}
