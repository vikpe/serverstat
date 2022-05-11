package qstatus_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/mvdsv/qstatus"
	"github.com/vikpe/serverstat/qserver/qsettings"
)

func TestParse(t *testing.T) {
	t.Run("Standby", func(t *testing.T) {
		settings := qsettings.Settings{
			"status":    "Standby",
			"timelimit": "10",
		}

		expect := qstatus.Status{
			Name:     "Standby",
			Duration: qstatus.MatchDuration{Elapsed: 0, Total: 10, Remaining: 10},
		}
		assert.Equal(t, expect, qstatus.Parse(settings))
	})

	t.Run("Started", func(t *testing.T) {
		settings := qsettings.Settings{
			"status":    "3 min left",
			"timelimit": "10",
		}

		expect := qstatus.Status{
			Name:     "Started",
			Duration: qstatus.MatchDuration{Elapsed: 7, Total: 10, Remaining: 3},
		}
		assert.Equal(t, expect, qstatus.Parse(settings))
	})
}
