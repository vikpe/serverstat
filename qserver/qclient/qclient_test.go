package qclient_test

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/mvdsv/commands/laststats"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qtext/qstring"
)

func TestNewFromString(t *testing.T) {
	t.Run("invalid", func(t *testing.T) {
		client, err := qclient.NewFromString("")
		assert.ErrorContains(t, err, "EOF")
		assert.Equal(t, client, qclient.Client{})
	})

	t.Run("missing fields", func(t *testing.T) {
		client, err := qclient.NewFromString("585 17 25")
		assert.ErrorContains(t, err, "invalid client column count 3, expects at least 8")
		assert.Equal(t, client, qclient.Client{})
	})

	t.Run("valid", func(t *testing.T) {
		expect := qclient.Client{
			Id:     585,
			Name:   qstring.New("XantoM"),
			Team:   qstring.New("f0m"),
			Skin:   "xantom",
			Colors: [2]uint8{4, 2},
			Frags:  17,
			Ping:   12,
			Time:   25,
			CC:     "SE",
		}
		clientString := `585 17 25 12 "XantoM" "xantom" 4 2 "f0m" "SE"`
		client, err := qclient.NewFromString(clientString)
		assert.Equal(t, expect, client)
		assert.Nil(t, err)
	})
}

func BenchmarkNewFromString(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		qclient.NewFromString(`585 17 25 12 "XantoM" "xantom" 4 2 "f0m" "SE"`)
	}
}

func TestFromStrings(t *testing.T) {
	clientStrings := []string{
		`63 5 4 25 "Pitbull" "" 4 4 "red"`,
		`66 2 4 38 "NL" "" 13 13 "red"`,
		`65 -9999 16 -666 "\s\[ServeMe]" "" 12 11 "lqwc"`,
		`67 -9999 122 -68 "\s\Final" "" 2 3 "red"`,
		``,
	}

	expect := []qclient.Client{
		{
			Id:     63,
			Name:   qstring.New("Pitbull"),
			Team:   qstring.New("red"),
			Skin:   "",
			Colors: [2]uint8{4, 4},
			Frags:  5,
			Ping:   25,
			Time:   4,
		},
		{
			Id:     66,
			Name:   qstring.New("NL"),
			Team:   qstring.New("red"),
			Skin:   "",
			Colors: [2]uint8{13, 13},
			Frags:  2,
			Ping:   38,
			Time:   4,
		},
		{
			Id:     65,
			Name:   qstring.New("[ServeMe]"),
			Team:   qstring.New("lqwc"),
			Skin:   "",
			Colors: [2]uint8{12, 11},
			Frags:  -9999,
			Ping:   -666,
			Time:   16,
		},
		{
			Id:     67,
			Name:   qstring.New("Final"),
			Team:   qstring.New("red"),
			Skin:   "",
			Colors: [2]uint8{2, 3},
			Frags:  -9999,
			Ping:   -68,
			Time:   122,
		},
	}

	actual := qclient.NewFromStrings(clientStrings)

	assert.Equal(t, expect, actual)
}

func TestNewFromLastStatsPlayer(t *testing.T) {
	lastStatsPlayer := laststats.Player{Name: "XantoM", Team: "f0m", Ping: 25, TopColor: 4, BottomColor: 2, Stats: laststats.PlayerStats{Frags: 6}}
	expect := qclient.Client{Name: "XantoM", Team: "f0m", Ping: 25, Colors: [2]uint8{4, 2}, Frags: 6}
	assert.Equal(t, expect, qclient.NewFromLastStatsPlayer(lastStatsPlayer))
}

func TestClient_IsSpectator(t *testing.T) {
	assert.True(t, qclient.Client{Ping: -10}.IsSpectator())
	assert.False(t, qclient.Client{Ping: 10}.IsSpectator())
}

func TestClient_IsPlayer(t *testing.T) {
	assert.False(t, qclient.Client{Ping: -10}.IsPlayer())
	assert.True(t, qclient.Client{Ping: 10}.IsPlayer())
}

func TestClient_IsBot(t *testing.T) {
	assert.True(t, qclient.Client{Name: qstring.New("XantoM"), Ping: 10}.IsBot())    // bot ping
	assert.True(t, qclient.Client{Name: qstring.New("[ServeMe]"), Ping: 12}.IsBot()) // bot name
	assert.False(t, qclient.Client{Name: qstring.New("XantoM"), Ping: 12}.IsBot())   // human name and ping
}

func TestClient_IsHuman(t *testing.T) {
	assert.False(t, qclient.Client{Name: qstring.New("XantoM"), Ping: 10}.IsHuman())    // bot ping
	assert.False(t, qclient.Client{Name: qstring.New("[ServeMe]"), Ping: 12}.IsHuman()) // bot name
	assert.True(t, qclient.Client{Name: qstring.New("XantoM"), Ping: 12}.IsHuman())     // human name and ping
}

func TestSortPlayers(t *testing.T) {
	milton := qclient.Client{Name: qstring.New("Milton"), Frags: 8}
	bps := qclient.Client{Name: qstring.New("bps"), Frags: 8}
	valla := qclient.Client{Name: qstring.New("valla"), Frags: 6}
	xantom := qclient.Client{Name: qstring.New("XantoM"), Frags: 12}
	xterm := qclient.Client{Name: qstring.New("Xterm"), Frags: 12}

	t.Run("no players", func(t *testing.T) {
		players := make([]qclient.Client, 0)
		expect := make([]qclient.Client, 0)
		assert.Equal(t, expect, qclient.SortPlayers(players))
	})

	t.Run("many players", func(t *testing.T) {
		players := []qclient.Client{milton, bps, xterm, valla, xantom}
		expect := []qclient.Client{xantom, xterm, bps, milton, valla}
		assert.Equal(t, expect, qclient.SortPlayers(players))
	})

	t.Run("few players", func(t *testing.T) {
		players := []qclient.Client{milton}
		expect := []qclient.Client{milton}
		assert.Equal(t, expect, qclient.SortPlayers(players))
	})
}

func BenchmarkSortPlayers(b *testing.B) {
	milton := qclient.Client{Name: qstring.New("Milton"), Frags: 8}
	bps := qclient.Client{Name: qstring.New("bps"), Frags: 8}
	valla := qclient.Client{Name: qstring.New("valla"), Frags: 6}
	xantom := qclient.Client{Name: qstring.New("XantoM"), Frags: 12}
	paradoks := qclient.Client{Name: qstring.New("ParadokS"), Frags: 21}
	players := []qclient.Client{milton, bps, valla, xantom, paradoks}

	b.ReportAllocs()
	b.ResetTimer()

	b.Run("no players", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			qclient.SortPlayers(make([]qclient.Client, 0))
		}
	})

	b.Run("one player", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			qclient.SortPlayers(players[0:1])
		}
	})

	b.Run("two players", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			qclient.SortPlayers(players[0:2])
		}
	})

	b.Run("many players", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			qclient.SortPlayers(players)
		}
	})
}

func TestClient_MarshalJSON(t *testing.T) {
	client := qclient.Client{
		Id:     67,
		Name:   qstring.New("Final"),
		Team:   qstring.New("red"),
		Skin:   "",
		Colors: [2]uint8{2, 3},
		Frags:  -9999,
		Ping:   -68,
		Time:   122,
	}

	jsonValue, _ := json.Marshal(client)
	expect := `{"id":67,"name":"Final","name_color":"wwwww","team":"red","team_color":"www","skin":"","colors":[2,3],"frags":-9999,"ping":-68,"time":122,"cc":"","is_bot":false}`
	assert.Equal(t, expect, string(jsonValue))
}

func TestClientNames(t *testing.T) {
	encodedNames := []string{
		"HCBtYXplcg==", // "• mazer",
		"EXNyEA==",     // "]sr[",
		"4uHz8w==",     // "bass",
	}

	names := make([]string, 0)

	for _, encodedString := range encodedNames {
		strBytes, _ := base64.StdEncoding.DecodeString(encodedString)
		names = append(names, string(strBytes))
	}

	clients := []qclient.Client{
		{Name: qstring.QuakeString(names[0])},
		{Name: qstring.QuakeString(names[1])},
		{Name: qstring.QuakeString(names[2])},
	}

	assert.Equal(t, []string{"• mazer", "]sr[", "bass"}, qclient.ClientNames(clients))
}
