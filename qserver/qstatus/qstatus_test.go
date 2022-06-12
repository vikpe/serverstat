package qstatus_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/qstatus"
)

func TestParse(t *testing.T) {
	type testCase struct {
		Status          string
		FreePlayerSlots int
		Expect          qstatus.Status
	}

	testCases := []testCase{
		{
			Status:          "Standby",
			FreePlayerSlots: 2,
			Expect: qstatus.Status{
				Name:        "Standby",
				Description: "Waiting for 2 players.",
			},
		},
		{
			Status:          "Standby",
			FreePlayerSlots: 1,
			Expect: qstatus.Status{
				Name:        "Standby",
				Description: "Waiting for 1 player.",
			},
		},
		{
			Status:          "Standby",
			FreePlayerSlots: 0,
			Expect: qstatus.Status{
				Name:        "Standby",
				Description: "Waiting for players to ready up.",
			},
		},
		{
			Status:          "Countdown",
			FreePlayerSlots: 0,
			Expect: qstatus.Status{
				Name:        "Started",
				Description: "Countdown",
			},
		},
		{
			Status:          "3 min left",
			FreePlayerSlots: 0,
			Expect: qstatus.Status{
				Name:        "Started",
				Description: "3 min left",
			},
		},
		{
			Status:          "foo",
			FreePlayerSlots: 2,
			Expect: qstatus.Status{
				Name:        "Unknown",
				Description: "foo",
			},
		},
	}

	for _, tc := range testCases {
		caseName := fmt.Sprintf("%s (%d free slots) = %s", tc.Status, tc.FreePlayerSlots, tc.Expect)

		t.Run(caseName, func(t *testing.T) {
			assert.Equal(t, tc.Expect, qstatus.New(tc.Status, tc.FreePlayerSlots), caseName)
		})
	}
}

func BenchmarkNew(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	b.Run("Standby", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			qstatus.New("Standby", 2)
		}
	})

	b.Run("x min left", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			qstatus.New("3 min left", 0)
		}
	})

	b.Run("Unknown", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			qstatus.New("foo", 4)
		}
	})
}
