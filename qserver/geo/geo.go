package geo

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/goccy/go-json"
	"github.com/vikpe/serverstat/qutil"
)

type Info struct {
	CC          string     `json:"cc"`
	Country     string     `json:"country"`
	Region      string     `json:"region"`
	City        string     `json:"city"`
	Coordinates [2]float32 `json:"coordinates"`
}

type IpToGeoMap map[string]Info

func NewIpToGeoMap(url string) IpToGeoMap {
	var result IpToGeoMap
	err := getJson(url, &result)

	if err != nil {
		timestamp := time.Now().Format(time.RFC850)
		fmt.Println(timestamp, "Unable to create geo ip map:", err.Error())
		return make(IpToGeoMap, 0)
	}

	return result
}

func (g IpToGeoMap) GetByAddress(address string) Info {
	hostname := strings.SplitN(address, ":", 2)[0]
	return g.GetByHostname(hostname)
}

func (g IpToGeoMap) GetByHostname(hostname string) Info {
	ip := qutil.HostnameToIp(hostname)
	return g.GetByIp(ip)
}

func (g IpToGeoMap) GetByIp(ip string) Info {
	if serverGeo, ok := g[ip]; ok {
		return serverGeo
	}

	return Info{
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
