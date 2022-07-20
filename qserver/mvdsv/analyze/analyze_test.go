package analyze_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/mvdsv/analyze"
	"github.com/vikpe/serverstat/qserver/mvdsv/qmode"
	"github.com/vikpe/serverstat/qserver/mvdsv/qstatus"
	"github.com/vikpe/serverstat/qserver/mvdsv/qtvstream"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qserver/qclient/slots"
	"github.com/vikpe/serverstat/qserver/qsettings"
	"github.com/vikpe/serverstat/qtext/qstring"
)

func TestGetPlayerNames(t *testing.T) {
	server := mvdsv.Mvdsv{
		Players: []qclient.Client{
			{Name: qstring.New("alpha")},
			{Name: qstring.New("beta")},
			{Name: qstring.New("gamma")},
		},
	}
	expect := []string{"alpha", "beta", "gamma"}
	assert.Equal(t, expect, analyze.GetPlayerPlainNames(server))
}

func TestHasPlayer(t *testing.T) {
	server := mvdsv.Mvdsv{
		Players: []qclient.Client{
			{Name: qstring.New("alpha")},
			{Name: qstring.New("beta")},
			{Name: qstring.New("gamma")},
		},
	}

	t.Run("no", func(t *testing.T) {
		assert.False(t, analyze.HasPlayer(server, "delta"))
	})

	t.Run("yes", func(t *testing.T) {
		assert.True(t, analyze.HasPlayer(server, "beta"))

		// wildcard matches
		assert.True(t, analyze.HasPlayer(server, "@eta"))
		assert.True(t, analyze.HasPlayer(server, "bet@"))
		assert.True(t, analyze.HasPlayer(server, "@et@"))
	})
}

func TestHasServerSpectator(t *testing.T) {
	server := mvdsv.Mvdsv{
		SpectatorNames: []string{"alpha", "beta", "gamma"},
	}

	t.Run("no", func(t *testing.T) {
		assert.False(t, analyze.HasServerSpectator(server, "delta"))
	})

	t.Run("yes", func(t *testing.T) {
		assert.True(t, analyze.HasServerSpectator(server, "beta"))

		// wildcard matches
		assert.True(t, analyze.HasServerSpectator(server, "@eta"))
		assert.True(t, analyze.HasServerSpectator(server, "bet@"))
		assert.True(t, analyze.HasServerSpectator(server, "@et@"))
	})
}

func TestHasQtvSpectator(t *testing.T) {
	server := mvdsv.Mvdsv{
		QtvStream: qtvstream.QtvStream{
			SpectatorNames: []string{"alpha", "beta", "gamma"},
		},
	}

	t.Run("no", func(t *testing.T) {
		assert.False(t, analyze.HasQtvSpectator(server, "delta"))
	})

	t.Run("yes", func(t *testing.T) {
		assert.True(t, analyze.HasQtvSpectator(server, "beta"))

		// wildcard matches
		assert.True(t, analyze.HasQtvSpectator(server, "@eta"))
		assert.True(t, analyze.HasQtvSpectator(server, "bet@"))
		assert.True(t, analyze.HasQtvSpectator(server, "@et@"))
	})
}

func TestHasSpectator(t *testing.T) {
	server := mvdsv.Mvdsv{
		SpectatorNames: []string{"alpha", "beta"},
		QtvStream: qtvstream.QtvStream{
			SpectatorNames: []string{"gamma"},
		},
	}

	t.Run("no", func(t *testing.T) {
		assert.False(t, analyze.HasSpectator(server, "delta"))
	})

	t.Run("yes", func(t *testing.T) {
		assert.True(t, analyze.HasSpectator(server, "beta"))
		assert.True(t, analyze.HasSpectator(server, "gamma"))
	})
}

func TestHasClient(t *testing.T) {
	server := mvdsv.Mvdsv{
		Players: []qclient.Client{
			{Name: qstring.New("alpha")},
			{Name: qstring.New("beta")},
		},
		SpectatorNames: []string{"gamma"},
		QtvStream: qtvstream.QtvStream{
			SpectatorNames: []string{"delta"},
		},
	}

	t.Run("no", func(t *testing.T) {
		assert.False(t, analyze.HasClient(server, "kappa"))
	})

	t.Run("yes", func(t *testing.T) {
		assert.True(t, analyze.HasClient(server, "beta"))
		assert.True(t, analyze.HasClient(server, "gamma"))
		assert.True(t, analyze.HasClient(server, "delta"))
	})
}

func TestIsIdle(t *testing.T) {
	t.Run("no players", func(t *testing.T) {
		server := mvdsv.Mvdsv{PlayerSlots: slots.New(4, 0)}
		assert.True(t, analyze.IsIdle(server))
	})

	t.Run("race", func(t *testing.T) {
		server := mvdsv.Mvdsv{
			PlayerSlots: slots.New(4, 1),
			Mode:        qmode.Mode("race"),
		}
		assert.False(t, analyze.IsIdle(server))
	})

	t.Run("2on2 started", func(t *testing.T) {
		server := mvdsv.Mvdsv{
			PlayerSlots: slots.New(4, 2),
			Mode:        qmode.Mode("2on2"),
			Status:      qstatus.Status{Name: "Started"},
		}
		assert.False(t, analyze.IsIdle(server))
	})

	t.Run("2on2 standby - not full server - players high time", func(t *testing.T) {
		server := mvdsv.Mvdsv{
			PlayerSlots: slots.New(4, 2),
			Mode:        qmode.Mode("2on2"),
			Status:      qstatus.Status{Name: "Standby"},
			Players:     []qclient.Client{{Time: 12}, {Time: 16}},
		}
		assert.False(t, analyze.IsIdle(server))
	})

	t.Run("1on1 standby - full server - players low time", func(t *testing.T) {
		server := mvdsv.Mvdsv{
			PlayerSlots: slots.New(2, 2),
			Mode:        qmode.Mode("1on1"),
			Status:      qstatus.Status{Name: "Standby"},
			Players:     []qclient.Client{{Time: 2}, {Time: 4}},
		}
		assert.False(t, analyze.IsIdle(server))
	})

	t.Run("1on1 standby - full server - players high time", func(t *testing.T) {
		server := mvdsv.Mvdsv{
			PlayerSlots: slots.New(2, 2),
			Mode:        qmode.Mode("1on1"),
			Status:      qstatus.Status{Name: "Standby"},
			Players:     []qclient.Client{{Time: 4}, {Time: 5}},
		}
		assert.True(t, analyze.IsIdle(server))
	})

	t.Run("2on2 standby - full server - players high time", func(t *testing.T) {
		server := mvdsv.Mvdsv{
			PlayerSlots: slots.New(4, 4),
			Mode:        qmode.Mode("2on2"),
			Status:      qstatus.Status{Name: "Standby"},
			Players:     []qclient.Client{{Time: 6}, {Time: 7}},
		}
		assert.True(t, analyze.IsIdle(server))
	})

	t.Run("coop standby - custom mode - player high time", func(t *testing.T) {
		server := mvdsv.Mvdsv{
			PlayerSlots: slots.New(26, 1),
			Mode:        qmode.Mode("coop"),
			Status:      qstatus.Status{Name: "Standby"},
			Players:     []qclient.Client{{Time: 3}},
		}
		assert.True(t, analyze.IsIdle(server))
	})
}

func TestMinPlayerTime(t *testing.T) {
	t.Run("no players", func(t *testing.T) {
		assert.Equal(t, 0, analyze.MinPlayerTime(nil))
	})

	t.Run("one players", func(t *testing.T) {
		players := []qclient.Client{{Time: 15}}
		assert.Equal(t, 15, analyze.MinPlayerTime(players))
	})

	t.Run("many players", func(t *testing.T) {
		players := []qclient.Client{{Time: 3}, {Time: 15}, {Time: 7}}
		assert.Equal(t, 3, analyze.MinPlayerTime(players))
	})
}

func TestRequiresPassword(t *testing.T) {
	testCases := map[string]bool{
		"0": false,
		"4": false,
		"5": false,
		"2": true,
		"3": true,
		"6": true,
		"7": true,
	}

	for needpass, expect := range testCases {

		t.Run(fmt.Sprintf("needpass=%s", needpass), func(t *testing.T) {
			server := mvdsv.Mvdsv{
				Settings: qsettings.Settings{
					"needpass": needpass,
				},
			}
			assert.Equal(t, expect, analyze.RequiresPassword(server))
		})
	}
}

func TestIsSpeccable(t *testing.T) {
	t.Run("yes - has qtv stream", func(t *testing.T) {
		server := mvdsv.Mvdsv{
			QtvStream: qtvstream.QtvStream{Url: "2@troopers.fi:28000"},
		}
		assert.True(t, analyze.IsSpeccable(server))
	})

	t.Run("yes - has free spectator slots and no password", func(t *testing.T) {
		server := mvdsv.Mvdsv{
			SpectatorSlots: slots.New(4, 3),
		}
		assert.True(t, analyze.IsSpeccable(server))
	})

	t.Run("no - has free spectator slots and password", func(t *testing.T) {
		server := mvdsv.Mvdsv{
			SpectatorSlots: slots.New(4, 3),
			Settings:       qsettings.Settings{"needpass": "2"},
		}
		assert.False(t, analyze.IsSpeccable(server))
	})
}
