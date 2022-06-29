package analyze_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/mvdsv/analyze"
	"github.com/vikpe/serverstat/qserver/mvdsv/qmode"
	"github.com/vikpe/serverstat/qserver/mvdsv/qstatus"
	"github.com/vikpe/serverstat/qserver/mvdsv/qtvstream"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qserver/qclient/slots"
	"github.com/vikpe/serverstat/qtext/qstring"
)

func TestListOfNamesContainsName(t *testing.T) {
	assert.False(t, analyze.ListOfNamesContainsName(nil, "foo"))
	assert.True(t, analyze.ListOfNamesContainsName([]string{"foo"}, "foo"))
}

func TestGetPlayerNames(t *testing.T) {
	server := mvdsv.Mvdsv{
		Players: []qclient.Client{
			{Name: qstring.New("alpha")},
			{Name: qstring.New("beta")},
			{Name: qstring.New("gamma")},
		},
	}
	expect := []string{"alpha", "beta", "gamma"}
	assert.Equal(t, expect, analyze.GetPlayerNames(server))
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
