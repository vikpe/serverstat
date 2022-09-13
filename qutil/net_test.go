package qutil_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qutil"
)

func TestHostnameToIp(t *testing.T) {
	assert.Equal(t, "91.102.91.59", qutil.HostnameToIp("qw.foppa.dk"))
	assert.Equal(t, "91.102.91.59", qutil.HostnameToIp("91.102.91.59"))
	assert.Equal(t, "foo", qutil.HostnameToIp("foo"))
}

func TestIsValidServerAddress(t *testing.T) {
	assert.False(t, qutil.IsValidServerAddress("foo"))
	assert.False(t, qutil.IsValidServerAddress("10.10.10.10"))

	assert.True(t, qutil.IsValidServerAddress("foo:28000"))
	assert.True(t, qutil.IsValidServerAddress("10.10.10.10:28000"))
}
