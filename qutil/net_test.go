package qutil_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qutil"
)

func TestHostnameToIp(t *testing.T) {
	assert.Equal(t, "91.102.91.59", qutil.HostnameToIp("qw.foppa.dk"))
	assert.Equal(t, "91.102.91.59", qutil.HostnameToIp("91.102.91.59"))
}
