package qutil_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qutil"
)

func TestClampInt(t *testing.T) {
	assert.Equal(t, 0, qutil.ClampInt(-5, 0, 10))
	assert.Equal(t, 10, qutil.ClampInt(15, 0, 10))
	assert.Equal(t, 5, qutil.ClampInt(5, 0, 10))
}
