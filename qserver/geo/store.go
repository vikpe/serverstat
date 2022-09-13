package geo

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/goccy/go-json"
	"github.com/vikpe/serverstat/qutil"
)

type Store map[string]Location

func NewStoreFromUrl(geoDataUrl string) Store {
	var result Store
	err := getJson(geoDataUrl, &result)

	if err != nil {
		timestamp := time.Now().Format(time.RFC850)
		fmt.Println(timestamp, "Unable to create geo store:", err.Error())
		return make(Store, 0)
	}

	return result
}

func (s Store) ByAddress(address string) Location {
	hostname := strings.SplitN(address, ":", 2)[0]
	return s.ByHostname(hostname)
}

func (s Store) ByHostname(hostname string) Location {
	ip := qutil.HostnameToIp(hostname)
	return s.ByIp(ip)
}

func (s Store) ByIp(ip string) Location {
	if info, ok := s[ip]; ok {
		return info
	}

	return Location{}
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
