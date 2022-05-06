package status32_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/mvdsv/commands/status32"
	"github.com/vikpe/serverstat/qserver/mvdsv/qtvstream"
)

func TestParseResponseBody(t *testing.T) {
	// empty response body
	result, err := status32.ParseResponse([]byte(""), nil)
	assert.Equal(t, errors.New("unable to parse response"), err)
	assert.Equal(t, qtvstream.QtvStream{}, result)

	// invalid qtv configuration
	result, err = status32.ParseResponse([]byte(`1 "qw.foppa.dk - qtv (3)" "" 2`), nil)
	assert.Equal(t, errors.New("invalid QTV configuration"), err)
	assert.Equal(t, qtvstream.QtvStream{}, result)

	// valid response body
	responseBody := []byte(`1 "qw.foppa.dk - qtv (3)" "3@qw.foppa.dk:28000" 4`)

	result, err = status32.ParseResponse(responseBody, nil)
	expect := qtvstream.QtvStream{
		Title:      "qw.foppa.dk - qtv (3)",
		Url:        "3@qw.foppa.dk:28000",
		NumClients: 4,
	}
	assert.Equal(t, nil, err)
	assert.Equal(t, expect, result)
}
