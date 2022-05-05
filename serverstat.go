package serverstat

import (
	"sync"

	"github.com/vikpe/serverstat/qserver"
)

func GetServerInfo(address string) (qserver.GenericServer, error) {
	return qserver.New(address)
}

func GetServerInfoFromMany(addresses []string) []qserver.GenericServer {
	var (
		wg    sync.WaitGroup
		mutex sync.Mutex
	)

	servers := make([]qserver.GenericServer, 0)

	for _, address := range addresses {
		wg.Add(1)

		go func(address string) {
			defer wg.Done()

			server, err := qserver.New(address)

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
