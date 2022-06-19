package qtime_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/qtime"
)

func TestNew(t *testing.T) {
	t.Run("Standby", func(t *testing.T) {
		expect := qtime.Time{
			Elapsed:   0,
			Total:     10,
			Remaining: 10,
		}
		time := qtime.New(10, "Standby")
		assert.Equal(t, expect, time)
	})

	t.Run("Countdown", func(t *testing.T) {
		expect := qtime.Time{
			Elapsed:   0,
			Total:     10,
			Remaining: 10,
		}
		time := qtime.New(10, "Countdown")
		assert.Equal(t, expect, time)
	})

	t.Run("Started", func(t *testing.T) {
		expect := qtime.Time{
			Elapsed:   7,
			Total:     10,
			Remaining: 3,
		}
		time := qtime.New(10, "3 min left")
		assert.Equal(t, expect, time)
	})
}
