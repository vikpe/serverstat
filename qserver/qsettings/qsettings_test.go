package qsettings_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/qsettings"
)

func TestParseString(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		assert.Equal(t, make(qsettings.Settings, 0), qsettings.ParseString(""))
	})
	t.Run("non-empty", func(t *testing.T) {
		expect := qsettings.Settings{
			"*version":  "MVDSV 0.35-dev",
			"maxfps":    "77",
			"pm_ktjump": "1",
		}
		settingsString := `maxfps\77\pm_ktjump\1\*version\MVDSV 0.35-dev`
		actual := qsettings.ParseString(settingsString)

		assert.Equal(t, expect, actual)
	})
}
