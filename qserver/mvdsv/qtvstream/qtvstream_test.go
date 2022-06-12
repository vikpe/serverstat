package qtvstream_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/mvdsv/qtvstream"
	"github.com/vikpe/serverstat/qtext/qstring"
)

func TestNew(t *testing.T) {
	expect := qtvstream.QtvStream{
		Title:          "",
		Url:            "",
		ID:             0,
		Address:        "",
		SpectatorNames: make([]qstring.QuakeString, 0),
		SpectatorCount: 0,
	}

	assert.Equal(t, expect, qtvstream.New())
}
