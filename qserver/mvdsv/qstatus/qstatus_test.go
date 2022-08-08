package qstatus_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/mvdsv/qmode"
	"github.com/vikpe/serverstat/qserver/mvdsv/qstatus"
	"github.com/vikpe/serverstat/qserver/qclient"
)

var botWithFrags = qclient.Client{Frags: 5, Ping: 10}
var humanWithFrags = qclient.Client{Frags: 7, Ping: 25}
var humanWithoutFrags = qclient.Client{Frags: 0, Ping: 25}

func TestParse(t *testing.T) {
	//botWithFrags = qclient.Client{Frags: 5, Ping: 10}
	//humanWithFrags = qclient.Client{Frags: 7, Ping: 25}

	type testCase struct {
		SettingsStatus string
		Mode           qmode.Mode
		Players        []qclient.Client
		FreeSlots      int
		Expect         qstatus.Status
	}

	testCases := []testCase{
		{
			SettingsStatus: "Standby",
			Mode:           "1on1",
			Players:        []qclient.Client{},
			FreeSlots:      2,
			Expect: qstatus.Status{
				Name:        "Standby",
				Description: "Waiting for 2 players",
			},
		},
		{
			SettingsStatus: "Standby",
			Mode:           "1on1",
			Players:        []qclient.Client{humanWithFrags, humanWithFrags},
			FreeSlots:      0,
			Expect: qstatus.Status{
				Name:        "Started",
				Description: "Score screen",
			},
		},
		{
			SettingsStatus: "Standby",
			Mode:           "1on1",
			Players:        []qclient.Client{humanWithFrags, botWithFrags},
			FreeSlots:      0,
			Expect: qstatus.Status{
				Name:        "Standby",
				Description: "Score screen",
			},
		},
		{
			SettingsStatus: "Standby",
			Mode:           "ffa",
			Players:        []qclient.Client{humanWithFrags, humanWithFrags},
			FreeSlots:      10,
			Expect: qstatus.Status{
				Name:        "Started",
				Description: "Score screen",
			},
		},
		{
			SettingsStatus: "Standby",
			Mode:           "ffa",
			Players:        []qclient.Client{humanWithFrags, botWithFrags},
			FreeSlots:      4,
			Expect: qstatus.Status{
				Name:        "Standby",
				Description: "Score screen",
			},
		},
		{
			SettingsStatus: "Standby",
			Mode:           "ffa",
			Players:        []qclient.Client{botWithFrags, botWithFrags},
			FreeSlots:      4,
			Expect: qstatus.Status{
				Name:        "Standby",
				Description: "Waiting for players",
			},
		},
		{
			SettingsStatus: "Standby",
			Mode:           "1on1",
			Players:        []qclient.Client{humanWithoutFrags},
			FreeSlots:      1,
			Expect: qstatus.Status{
				Name:        "Standby",
				Description: "Waiting for 1 player",
			},
		},
		{
			SettingsStatus: "Standby",
			Mode:           "1on1",
			Players:        []qclient.Client{humanWithoutFrags, humanWithoutFrags},
			FreeSlots:      0,
			Expect: qstatus.Status{
				Name:        "Standby",
				Description: "Waiting for players to ready up",
			},
		},
		{
			SettingsStatus: "Countdown",
			Mode:           "1on1",
			Players:        []qclient.Client{humanWithoutFrags, humanWithoutFrags},
			FreeSlots:      0,
			Expect: qstatus.Status{
				Name:        "Started",
				Description: "Countdown",
			},
		},
		{
			SettingsStatus: "3 min left",
			Mode:           "1on1",
			Players:        []qclient.Client{humanWithFrags, humanWithFrags},
			FreeSlots:      0,
			Expect: qstatus.Status{
				Name:        "Started",
				Description: "3 min left",
			},
		},
		{
			SettingsStatus: "Standby",
			Mode:           "coop",
			Players:        []qclient.Client{humanWithoutFrags, humanWithoutFrags, humanWithoutFrags},
			FreeSlots:      4,
			Expect: qstatus.Status{
				Name:        "Standby",
				Description: "Waiting for players to ready up",
			},
		},
		{
			SettingsStatus: "0 min left",
			Mode:           "coop",
			Players:        []qclient.Client{humanWithFrags},
			FreeSlots:      4,
			Expect: qstatus.Status{
				Name:        "Started",
				Description: "Game in progress",
			},
		},
		{
			SettingsStatus: "Standby",
			Mode:           "race",
			Players:        []qclient.Client{humanWithoutFrags},
			FreeSlots:      4,
			Expect: qstatus.Status{
				Name:        "Standby",
				Description: "Racing",
			},
		},
		{
			SettingsStatus: "foo",
			Mode:           "coop",
			Players:        []qclient.Client{humanWithoutFrags, humanWithoutFrags, humanWithoutFrags, humanWithoutFrags},
			FreeSlots:      0,
			Expect: qstatus.Status{
				Name:        "Unknown",
				Description: "foo",
			},
		},
		{
			SettingsStatus: "Standby",
			Mode:           "clan arena",
			Players:        []qclient.Client{humanWithoutFrags, humanWithoutFrags, humanWithoutFrags, humanWithoutFrags},
			FreeSlots:      16,
			Expect: qstatus.Status{
				Name:        "Standby",
				Description: "Waiting for players to ready up",
			},
		},
		{
			SettingsStatus: "27 min left",
			Mode:           "clan arena",
			Players:        []qclient.Client{humanWithoutFrags, humanWithoutFrags, humanWithoutFrags, humanWithoutFrags},
			FreeSlots:      16,
			Expect: qstatus.Status{
				Name:        "Started",
				Description: "Game in progress",
			},
		},
	}

	for _, tc := range testCases {
		caseName := fmt.Sprintf("%s %s (%d free slots) = %s", tc.Mode, tc.SettingsStatus, tc.FreeSlots, tc.Expect)

		t.Run(caseName, func(t *testing.T) {
			assert.Equal(t, tc.Expect, qstatus.New(tc.SettingsStatus, tc.Mode, tc.Players, tc.FreeSlots), caseName)
		})
	}
}

func BenchmarkNew(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	b.Run("Standby - waiting", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			qstatus.New("Standby", "1on1", []qclient.Client{humanWithoutFrags, humanWithoutFrags}, 2)
		}
	})

	b.Run("Started - x min left", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			qstatus.New("3 min left", "1on1", []qclient.Client{botWithFrags, humanWithFrags}, 0)
		}
	})

	b.Run("Unknown", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			qstatus.New("foo", "coop", []qclient.Client{humanWithoutFrags, humanWithoutFrags, humanWithoutFrags, humanWithoutFrags}, 6)
		}
	})
}
