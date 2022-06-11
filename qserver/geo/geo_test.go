package geo_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/geo"
)

func TestIpGeoMap(t *testing.T) {
	ipToGeo := geo.NewIpToGeoMap()
	assert.True(t, len(ipToGeo) > 100)

	// test methods
	var ip string
	var info geo.Info

	for ip_, info_ := range ipToGeo {
		ip = ip_
		info = info_
		break
	}

	t.Run("GetByIp", func(t *testing.T) {
		t.Run("known", func(t *testing.T) {
			assert.Equal(t, info, ipToGeo.GetByIp(ip))

		})

		t.Run("unknown", func(t *testing.T) {
			expect := geo.Info{
				CC:          "",
				Country:     "",
				Region:      "",
				City:        "",
				Coordinates: [2]float32{0, 0},
			}
			assert.Equal(t, expect, ipToGeo.GetByIp("zzz"))
		})
	})

	t.Run("GetByAddress", func(t *testing.T) {
		assert.Equal(t, info, ipToGeo.GetByAddress(fmt.Sprintf("%s:%d", ip, 28000)))
	})
}
