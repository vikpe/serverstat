package convert

import (
	"github.com/ssoroka/slice"
	"github.com/vikpe/serverstat/qserver"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qserver/qtv"
	"github.com/vikpe/serverstat/qserver/qwfwd"
	"github.com/vikpe/serverstat/qtext/qstring"
)

func ToMvdsv(server qserver.GenericServer) mvdsv.Mvdsv {
	spectatorNames := slice.Map[qclient.Client, qstring.QuakeString](server.Spectators(), func(client qclient.Client) qstring.QuakeString {
		return client.Name
	})

	return mvdsv.Mvdsv{
		Address:        server.Address,
		Players:        server.Players(),
		SpectatorNames: spectatorNames,
		Settings:       server.Settings,
		QtvStream:      server.ExtraInfo.QtvStream,
		Geo:            server.ExtraInfo.Geo,
	}
}

func ToMvdsvExport(server qserver.GenericServer) mvdsv.MvdsvExport {
	return mvdsv.Export(ToMvdsv(server))
}

func ToQtv(server qserver.GenericServer) qtv.Qtv {
	spectatorNames := make([]qstring.QuakeString, 0)

	for _, client := range server.Clients {
		spectatorNames = append(spectatorNames, client.Name)
	}

	return qtv.Qtv{
		Address:        server.Address,
		SpectatorNames: spectatorNames,
		Settings:       server.Settings,
		Geo:            server.ExtraInfo.Geo,
	}
}

func ToQtvExport(server qserver.GenericServer) qtv.QtvExport {
	return qtv.Export(ToQtv(server))
}

func ToQwfwd(server qserver.GenericServer) qwfwd.Qwfwd {
	clientNames := slice.Map[qclient.Client, qstring.QuakeString](server.Clients, func(client qclient.Client) qstring.QuakeString {
		return client.Name
	})

	return qwfwd.Qwfwd{
		Address:     server.Address,
		ClientNames: clientNames,
		Settings:    server.Settings,
		Geo:         server.ExtraInfo.Geo,
	}
}

func ToQwfwdExport(server qserver.GenericServer) qwfwd.QwfwdExport {
	return qwfwd.Export(ToQwfwd(server))
}
