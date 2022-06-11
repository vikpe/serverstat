package geo

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/goccy/go-json"
)

type Info struct {
	CC          string
	Country     string
	Region      string
	City        string
	Coordinates [2]float32
}

type IpToGeoMap map[string]Info

func NewIpToGeoMap() IpToGeoMap {
	url := "https://raw.githubusercontent.com/vikpe/qw-servers-geoip/main/ip_to_geo.json"

	var result IpToGeoMap
	err := getJson(url, &result)

	if err != nil {
		fmt.Println("Unable to create geo ip map")
		return make(IpToGeoMap, 0)
	}

	return result
}

func (g IpToGeoMap) GetByAddress(address string) Info {
	ip := strings.SplitN(address, ":", 2)[0]
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
	httpClient := &http.Client{Timeout: 5 * time.Second}
	response, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return json.NewDecoder(response.Body).Decode(target)
}
