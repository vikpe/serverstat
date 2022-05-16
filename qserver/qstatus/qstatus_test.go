package qstatus_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/qstatus"
)

func TestParse(t *testing.T) {
	testCases := map[string]string{
		"Standby":    "Standby",
		"Countdown":  "Countdown",
		"Started":    "Started",
		"3 min left": "Started",
	}

	for status, expect := range testCases {
		t.Run(status, func(t *testing.T) {
			assert.Equal(t, expect, qstatus.Parse(status))
		})
	}
}
