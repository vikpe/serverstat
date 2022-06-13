package serverstat

import (
	"sort"
	"sync"

	"github.com/vikpe/serverstat/qserver"
	"github.com/vikpe/serverstat/qserver/commands/status87"
	"github.com/vikpe/serverstat/qserver/geo"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/mvdsv/qtvstream"
	"github.com/vikpe/serverstat/qserver/qversion"
	"github.com/vikpe/udpclient"
)

func GetInfo(address string) (qserver.GenericServer, error) {
	server, err := getServerInfo(address)

	if err == nil {
		ipToGeoMap := geo.NewIpToGeoMap()
		server.Geo = ipToGeoMap.GetByAddress(server.Address)
	}

	return server, err
}

func getServerInfo(address string) (qserver.GenericServer, error) {
	settings, clients, err := status87.ParseResponse(
		udpclient.New().SendCommand(address, status87.Command),
	)

	if err != nil {
		return qserver.GenericServer{}, err
	}

	if settings.Has("hostname") {
		settings["hostname_parsed"] = qserver.ParseHostname(address, settings.Get("hostname", ""))
	}

	server := qserver.GenericServer{
		Address:  address,
		Version:  qversion.New(settings.Get("*version", "")),
		Clients:  clients,
		Settings: settings,
	}

	if server.Version.IsMvdsv() {
		stream, _ := mvdsv.GetQtvStream(address)
		server.ExtraInfo.QtvStream = stream
	} else {
		server.ExtraInfo.QtvStream = qtvstream.New()
	}

	return server, nil
}

func GetInfoFromMany(addresses []string) []qserver.GenericServer {
	var (
		wg    sync.WaitGroup
		mutex sync.Mutex
	)

	servers := make([]qserver.GenericServer, 0)

	for _, address := range addresses {
		wg.Add(1)

		go func(address string) {
			defer wg.Done()

			server, err := getServerInfo(address)

			if err != nil {
				return
			}

			mutex.Lock()
			servers = append(servers, server)
			mutex.Unlock()
		}(address)
	}

	wg.Wait()

	ipToGeo := geo.NewIpToGeoMap()

	for index, server := range servers {
		servers[index].Geo = ipToGeo.GetByAddress(server.Address)
	}

	sort.Slice(servers, func(i, j int) bool {
		return servers[i].Address < servers[j].Address
	})

	return servers
}
