package qutil_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qutil"
)

func TestStringToInt(t *testing.T) {
	assert.Equal(t, 1, qutil.StringToInt("1"))
	assert.Equal(t, 55, qutil.StringToInt("55"))
	assert.Equal(t, -55, qutil.StringToInt("-55"))
	assert.Equal(t, 0, qutil.StringToInt("0"))
	assert.Equal(t, 0, qutil.StringToInt(""))
}
