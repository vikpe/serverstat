package qserver_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qtext/qstring"
)

func TestGenericServer_Players(t *testing.T) {
	player1 := qclient.Client{Name: qstring.New("player 1"), Ping: 25}
	player2 := qclient.Client{Name: qstring.New("player 2"), Ping: 12}
	spectator1 := qclient.Client{Name: qstring.New("spectator"), Ping: -25}
	server := qserver.GenericServer{
		Clients: []qclient.Client{player1, player2, spectator1},
	}

	expect := []qclient.Client{player1, player2}
	assert.Equal(t, expect, server.Players())
}

func TestGenericServer_Spectators(t *testing.T) {
	player1 := qclient.Client{Name: qstring.New("player 1"), Ping: 25}
	player2 := qclient.Client{Name: qstring.New("player 2"), Ping: 12}
	spectator1 := qclient.Client{Name: qstring.New("spectator"), Ping: -25}
	server := qserver.GenericServer{
		Clients: []qclient.Client{player1, player2, spectator1},
	}

	expect := []qclient.Client{spectator1}
	assert.Equal(t, expect, server.Spectators())
}

func TestParseAddress(t *testing.T) {
	testCases := map[string][2]string{
		"10.10.10.10:28501":     {"10.10.10.10:28501", ""},
		"10.10.10.10:28502":     {"10.10.10.10:28502", "foo:bar"},
		"78.141.238.193:27500":  {"78.141.238.193:27500", "Barrysworld FFA Tribute"},
		"91.206.14.17:27508":    {"91.206.14.17:27508", "Gladius SPB KTX #27508"},
		"108.61.178.207:28505":  {"108.61.178.207:28505", "de.aye.wtf:2850"},
		"213.239.216.253:27500": {"213.239.216.253:27500", "FFA @ qw.servegame.org�"},
		"91.102.91.59:27501":    {"91.102.91.59:27501", "qw.foppa.dk#1 - ktx"},
		"qw.foppa.dk:27501":     {"91.102.91.59:27501", "qw.foppa.dk #1 - ktx"},
		"de.aye.wtf:28506":      {"108.61.178.207:28506", "de.aye.wtf:28506"},
		"de.aye.wtf:28503":      {"108.61.178.207:28503", "de.aye.wtf:28503�"},
		"qw.irc.ax:28501":       {"46.227.68.148:28501", "QW.IRC.AX KTX:28501 (1 vs. trl)�"},
	}

	for expect, input := range testCases {
		assert.Equal(t, expect, qserver.ParseAddress(input[0], input[1]))
	}
}

func BenchmarkParseAddress(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	b.Run("no hostname", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			qserver.ParseAddress("78.141.238.193:27500", "")
		}
	})

	b.Run("invalid hostname", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			qserver.ParseAddress("78.141.238.193:27500", "Barrysworld FFA Tribute")
		}
	})

	b.Run("partial hostname", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			qserver.ParseAddress("91.102.91.59:27501", "qw.foppa.dk #1 - ktx")
		}
	})

	b.Run("exact hostname", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			qserver.ParseAddress("108.61.178.207:28506", "de.aye.wtf:28506")
		}
	})
}
