package convert_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver"
	"github.com/vikpe/serverstat/qserver/convert"
	"github.com/vikpe/serverstat/qserver/geo"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/mvdsv/qmode"
	"github.com/vikpe/serverstat/qserver/mvdsv/qstatus"
	"github.com/vikpe/serverstat/qserver/mvdsv/qtvstream"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qserver/qclient/slots"
	"github.com/vikpe/serverstat/qserver/qsettings"
	"github.com/vikpe/serverstat/qserver/qteam"
	"github.com/vikpe/serverstat/qserver/qtime"
	"github.com/vikpe/serverstat/qserver/qtv"
	"github.com/vikpe/serverstat/qserver/qwfwd"
	"github.com/vikpe/serverstat/qtext/qstring"
)

var PlayerClient = qclient.Client{
	Name:   qstring.New("XantoM"),
	Team:   qstring.New("red"),
	Skin:   "",
	Colors: [2]uint8{13, 13},
	Frags:  2,
	Ping:   38,
	Time:   4,
}

var SpectatorClient = qclient.Client{
	Name:   qstring.New("[ServeMe]"),
	Team:   qstring.New("lqwc"),
	Skin:   "",
	Colors: [2]uint8{12, 11},
	Frags:  -9999,
	Ping:   -666,
	Time:   16,
}

var GenericServer = qserver.GenericServer{
	Address:  "qw.foppa.dk:27501",
	Version:  "mvdsv 0.15",
	Clients:  []qclient.Client{PlayerClient, SpectatorClient},
	Settings: qsettings.Settings{"map": "dm2", "*gamedir": "qw", "status": "3 min left", "timelimit": "10", "maxclients": "8", "maxspectators": "4", "teamplay": "2"},
	Geo:      geo.Location{},
	ExtraInfo: struct {
		QtvStream qtvstream.QtvStream `json:"qtv_stream"`
	}{},
}

func TestToMvdsv(t *testing.T) {
	expect := mvdsv.Mvdsv{
		Address: GenericServer.Address,
		Mode:    qmode.Mode("4on4"),
		Submode: "",
		Title:   "4on4: red (XantoM) [dm2]",
		Status: qstatus.Status{
			Name:        "Started",
			Description: "3 min left",
		},
		Time: qtime.Time{
			Elapsed:   7,
			Total:     10,
			Remaining: 3,
		},
		Players: []qclient.Client{PlayerClient},
		PlayerSlots: slots.Slots{
			Used:  1,
			Total: 8,
			Free:  7,
		},
		SpectatorNames: []string{"[ServeMe]"},
		SpectatorSlots: slots.Slots{
			Used:  1,
			Total: 4,
			Free:  3,
		},
		Settings:  GenericServer.Settings,
		QtvStream: GenericServer.ExtraInfo.QtvStream,
		Teams: []qteam.Team{
			{
				Name:    qstring.New("red"),
				Players: []qclient.Client{PlayerClient},
			},
		},
		Geo: geo.Location{
			CC:      "",
			Country: "",
			Region:  "",
		},
		Score: 5,
	}

	assert.Equal(t, expect, convert.ToMvdsv(GenericServer))
}

func BenchmarkToMvdsv(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		convert.ToMvdsv(GenericServer)
	}
}

func TestToQtv(t *testing.T) {
	expect := qtv.Qtv{
		Address:        GenericServer.Address,
		SpectatorNames: []string{"XantoM", "[ServeMe]"},
		Settings:       GenericServer.Settings,
	}

	assert.Equal(t, expect, convert.ToQtv(GenericServer))
}

func BenchmarkToQtv(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		convert.ToQtv(GenericServer)
	}
}

func TestToQwfwd(t *testing.T) {
	expect := qwfwd.Qwfwd{
		Address:     GenericServer.Address,
		ClientNames: []string{"XantoM", "[ServeMe]"},
		Settings:    GenericServer.Settings,
	}

	assert.Equal(t, expect, convert.ToQwfwd(GenericServer))
}

func TestToJson(t *testing.T) {
	expect := `{"address":"qw.foppa.dk:27501","mode":"4on4","submode":"","title":"4on4: red (XantoM) [dm2]","status":{"name":"Started","description":"3 min left"},"time":{"elapsed":7,"total":10,"remaining":3},"player_slots":{"used":1,"total":8,"free":7},"players":[{"name":"XantoM","name_color":"wwwwww","team":"red","team_color":"www","skin":"","colors":[13,13],"frags":2,"ping":38,"time":4,"cc":"","is_bot":false}],"teams":[{"name":"red","name_color":"www","frags":2,"ping":38,"colors":[13,13],"players":[{"name":"XantoM","name_color":"wwwwww","team":"red","team_color":"www","skin":"","colors":[13,13],"frags":2,"ping":38,"time":4,"cc":"","is_bot":false}]}],"spectator_slots":{"used":1,"total":4,"free":3},"spectator_names":["[ServeMe]"],"settings":{"*gamedir":"qw","map":"dm2","maxclients":"8","maxspectators":"4","status":"3 min left","teamplay":"2","timelimit":"10"},"qtv_stream":{"title":"","url":"","id":0,"address":"","spectator_names":null,"spectator_count":0},"geo":{"cc":"","country":"","region":"","city":"","coordinates":[0,0]},"score":5,"type":"mvdsv"}`

	assert.JSONEq(t, expect, convert.ToJson(GenericServer))
}
