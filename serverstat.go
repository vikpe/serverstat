package serverstat

import (
	"sync"

	"github.com/vikpe/serverstat/qserver"
	"github.com/vikpe/serverstat/qserver/commands/status23"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/udpclient"
)

func GetInfo(address string) (qserver.GenericServer, error) {
	server, err := status23.Send(udpclient.New(), address)

	if server.Version.IsMvdsv() {
		server.ExtraInfo.QtvStream, _ = mvdsv.GetQtvStream(address)
	}

	return server, err
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

			server, err := GetInfo(address)

			if err != nil {
				return
			}

			mutex.Lock()
			servers = append(servers, server)
			mutex.Unlock()
		}(address)
	}

	wg.Wait()

	return servers
}
