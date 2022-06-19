package analyze_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/mvdsv/analyze"
	"github.com/vikpe/serverstat/qserver/mvdsv/qmode"
	"github.com/vikpe/serverstat/qserver/mvdsv/qstatus"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qserver/qclient/slots"
)

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

	t.Run("not standby", func(t *testing.T) {
		server := mvdsv.Mvdsv{
			PlayerSlots: slots.New(4, 2),
			Mode:        qmode.Mode("2on2"),
			Status:      qstatus.Status{Name: "Started"},
		}
		assert.False(t, analyze.IsIdle(server))
	})

	t.Run("standby - not full server - players high time", func(t *testing.T) {
		server := mvdsv.Mvdsv{
			PlayerSlots: slots.New(4, 2),
			Mode:        qmode.Mode("2on2"),
			Status:      qstatus.Status{Name: "Standby"},
			Players:     []qclient.Client{{Time: 12}, {Time: 16}},
		}
		assert.False(t, analyze.IsIdle(server))
	})

	t.Run("standby - full server - players low time", func(t *testing.T) {
		server := mvdsv.Mvdsv{
			PlayerSlots: slots.New(2, 2),
			Mode:        qmode.Mode("1on1"),
			Status:      qstatus.Status{Name: "Standby"},
			Players:     []qclient.Client{{Time: 3}, {Time: 9}},
		}
		assert.False(t, analyze.IsIdle(server))
	})

	t.Run("standby - full server - players high time", func(t *testing.T) {
		server := mvdsv.Mvdsv{
			PlayerSlots: slots.New(2, 2),
			Mode:        qmode.Mode("1on1"),
			Status:      qstatus.Status{Name: "Standby"},
			Players:     []qclient.Client{{Time: 15}, {Time: 32}},
		}
		assert.True(t, analyze.IsIdle(server))
	})

	t.Run("standby - custom mode - player high time", func(t *testing.T) {
		server := mvdsv.Mvdsv{
			PlayerSlots: slots.New(26, 1),
			Mode:        qmode.Mode("coop"),
			Status:      qstatus.Status{Name: "Standby"},
			Players:     []qclient.Client{{Time: 32}},
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
