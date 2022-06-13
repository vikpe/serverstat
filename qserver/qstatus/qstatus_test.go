package qstatus_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/mvdsv/qmode"
	"github.com/vikpe/serverstat/qserver/qstatus"
)

func TestParse(t *testing.T) {
	type testCase struct {
		SettingsStatus  string
		FreePlayerSlots int
		Mode            qmode.Mode
		Expect          qstatus.Status
	}

	testCases := []testCase{
		{
			SettingsStatus:  "Standby",
			FreePlayerSlots: 2,
			Mode:            "1on1",
			Expect: qstatus.Status{
				Name:        "Standby",
				Description: "Waiting for 2 players",
			},
		},
		{
			SettingsStatus:  "Standby",
			FreePlayerSlots: 1,
			Mode:            "1on1",
			Expect: qstatus.Status{
				Name:        "Standby",
				Description: "Waiting for 1 player",
			},
		},
		{
			SettingsStatus:  "Standby",
			FreePlayerSlots: 0,
			Mode:            "1on1",
			Expect: qstatus.Status{
				Name:        "Standby",
				Description: "Waiting for players to ready up",
			},
		},
		{
			SettingsStatus:  "Countdown",
			FreePlayerSlots: 0,
			Mode:            "1on1",
			Expect: qstatus.Status{
				Name:        "Started",
				Description: "Countdown",
			},
		},
		{
			SettingsStatus:  "3 min left",
			FreePlayerSlots: 0,
			Mode:            "1on1",
			Expect: qstatus.Status{
				Name:        "Started",
				Description: "3 min left",
			},
		},
		{
			SettingsStatus:  "Standby",
			FreePlayerSlots: 8,
			Mode:            "coop",
			Expect: qstatus.Status{
				Name:        "Standby",
				Description: "Waiting for players to ready up",
			},
		},
		{
			SettingsStatus:  "0 min left",
			FreePlayerSlots: 5,
			Mode:            "coop",
			Expect: qstatus.Status{
				Name:        "Started",
				Description: "Game in progress",
			},
		},
		{
			SettingsStatus:  "foo",
			FreePlayerSlots: 2,
			Mode:            "1on1",
			Expect: qstatus.Status{
				Name:        "Unknown",
				Description: "foo",
			},
		},
	}

	for _, tc := range testCases {
		caseName := fmt.Sprintf("%s (%d free slots) = %s", tc.SettingsStatus, tc.FreePlayerSlots, tc.Expect)

		t.Run(caseName, func(t *testing.T) {
			assert.Equal(t, tc.Expect, qstatus.New(tc.SettingsStatus, tc.FreePlayerSlots, tc.Mode), caseName)
		})
	}
}

func BenchmarkNew(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	b.Run("Standby", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			qstatus.New("Standby", 2, "1on1")
		}
	})

	b.Run("x min left", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			qstatus.New("3 min left", 0, "1on1")
		}
	})

	b.Run("Unknown", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			qstatus.New("foo", 4, "1on1")
		}
	})
}
