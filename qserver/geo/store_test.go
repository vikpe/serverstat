package geo_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/geo"
)

func TestNewGeoStoreFromUrl(t *testing.T) {
	store := geo.NewStoreFromUrl("https://raw.githubusercontent.com/vikpe/qw-servers-geoip/main/ip_to_geo.json")
	assert.True(t, len(store) > 100)

	// test methods
	var ip string
	var info geo.Location

	for ip_, info_ := range store {
		ip = ip_
		info = info_
		break
	}

	t.Run("ByIp", func(t *testing.T) {
		t.Run("known", func(t *testing.T) {
			assert.Equal(t, info, store.ByIp(ip))
		})

		t.Run("unknown", func(t *testing.T) {
			expect := geo.Location{
				CC:          "",
				Country:     "",
				Region:      "",
				City:        "",
				Coordinates: [2]float32{0, 0},
			}
			assert.Equal(t, expect, store.ByIp("zzz"))
		})
	})

	t.Run("ByAddress", func(t *testing.T) {
		assert.Equal(t, info, store.ByAddress(fmt.Sprintf("%s:%d", ip, 28000)))
	})
}
