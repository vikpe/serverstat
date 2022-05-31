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
