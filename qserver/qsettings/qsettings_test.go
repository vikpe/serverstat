package qsettings_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/qsettings"
)

func TestParseString(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		assert.Equal(t, qsettings.Settings{}, qsettings.ParseString(""))
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

func TestSettings_Has(t *testing.T) {
	settings := qsettings.Settings{"foo": "bar"}
	assert.True(t, settings.Has("foo"))
	assert.False(t, settings.Has("alpha"))
}

func TestSettings_Get(t *testing.T) {
	settings := qsettings.Settings{"foo": "bar"}
	assert.Equal(t, "bar", settings.Get("foo", ""))
	assert.Equal(t, "", settings.Get("alpha", ""))
}

func TestSettings_GetInt(t *testing.T) {
	settings := qsettings.Settings{"foo": "2"}
	assert.Equal(t, 2, settings.GetInt("foo", 0))
	assert.Equal(t, 0, settings.GetInt("alpha", 0))
}
