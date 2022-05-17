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

func TestReverse(t *testing.T) {
	assert.Equal(t, "321ateb", qutil.ReverseString("beta123"))
}

func TestCommonPrefix(t *testing.T) {
	assert.Equal(t, "", qutil.CommonPrefix([]string{"alpha", "beta"}))
	assert.Equal(t, "foo-", qutil.CommonPrefix([]string{"foo-alpha", "foo-beta"}))
}

func TestCommonSuffix(t *testing.T) {
	assert.Equal(t, "", qutil.CommonSuffix([]string{"foo", "bar"}))
	assert.Equal(t, "a-foo", qutil.CommonSuffix([]string{"alpha-foo", "beta-foo"}))
}
