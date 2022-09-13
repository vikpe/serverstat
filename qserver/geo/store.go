package geo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/vikpe/gonet"
)

type Store map[string]Location

func NewStoreFromUrl(geoDataUrl string) Store {
	var result Store
	err := getJson(geoDataUrl, &result)

	if err != nil {
		timestamp := time.Now().Format(time.RFC850)
		fmt.Println(timestamp, "Unable to create geo ip map:", err.Error())
		return make(Store, 0)
	}

	return result
}

func (g Store) ByAddress(address string) Location {
	hostname := strings.SplitN(address, ":", 2)[0]
	return g.ByHostname(hostname)
}

func (g Store) ByHostname(hostname string) Location {
	ip, _ := gonet.ToIpV4(hostname)
	return g.ByIp(ip)
}

func (g Store) ByIp(ip string) Location {
	if info, ok := g[ip]; ok {
		return info
	}

	return Location{
		CC:          "",
		Country:     "",
		Region:      "",
		City:        "",
		Coordinates: [2]float32{0, 0},
	}
}

func getJson(url string, target interface{}) error {
	httpClient := &http.Client{Timeout: 3 * time.Second}
	response, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return json.NewDecoder(response.Body).Decode(target)
}
