package qsettings_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/qsettings"
)

func TestNew(t *testing.T) {
	// empty
	assert.Equal(t, make(map[string]string, 0), qsettings.New(""))

	// non-empty
	expect := map[string]string{
		"*version":  "MVDSV 0.35-dev",
		"maxfps":    "77",
		"pm_ktjump": "1",
	}
	settingsString := `maxfps\77\pm_ktjump\1\*version\MVDSV 0.35-dev`
	actual := qsettings.New(settingsString)

	assert.Equal(t, expect, actual)
}
