package qtvusers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/mvdsv/commands/qtvusers"
	"github.com/vikpe/serverstat/qserver/qclient"
)

func TestParseResponseBody(t *testing.T) {
	// empty response body
	expect := make([]qclient.Client, 0)
	result := qtvusers.ParseResponseBody([]byte("\n"))
	assert.Equal(t, expect, result)

	// non-empty response body
	responseBody := []byte(`12 "djevulsk" "serp"`)

	result = qtvusers.ParseResponseBody(responseBody)
	expect = []qclient.Client{
		{
			Name:    "djevulsk",
			NameRaw: []int32{100, 106, 101, 118, 117, 108, 115, 107},
		},
		{
			Name:    "serp",
			NameRaw: []int32{115, 101, 114, 112},
		},
	}

	assert.Equal(t, expect, result)
}
