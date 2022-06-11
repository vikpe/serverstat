package geo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/geo"
)

func TestNewIpToGeoMap(t *testing.T) {
	ipToGeo := geo.NewIpToGeoMap()

	assert.True(t, len(ipToGeo) > 100)

	// test first entry
	var ip string
	var info geo.Info

	for ip_, info_ := range ipToGeo {
		ip = ip_
		info = info_
		break
	}

	assert.Equal(t, info, ipToGeo.GetByIp(ip))
}
