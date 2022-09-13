package geo_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/geo"
)

var ip = "100.11.123.208"
var location = geo.Location{
	CC:          "US",
	Country:     "United States",
	Region:      "North America",
	City:        "Lansdale",
	Coordinates: [2]float32{40.2363, -75.296},
}
var store = geo.Store{ip: location}
var nullLocation = geo.Location{}

func TestStore_ByIp(t *testing.T) {
	t.Run("known", func(t *testing.T) {
		assert.Equal(t, location, store.ByIp(ip))
	})

	t.Run("unknown", func(t *testing.T) {
		assert.Equal(t, nullLocation, store.ByIp("1.1.1.1"))
	})
}

func TestStore_ByAddress(t *testing.T) {
	t.Run("known", func(t *testing.T) {
		assert.Equal(t, location, store.ByAddress(fmt.Sprintf("%s:28000", ip)))
	})

	t.Run("unknown", func(t *testing.T) {
		assert.Equal(t, nullLocation, store.ByAddress("foo:28000"))
	})
}

func TestNewStoreFromUrl(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		storeFromUrl := geo.NewStoreFromUrl("foo")
		assert.Equal(t, geo.Store{}, storeFromUrl)
	})

	t.Run("success", func(t *testing.T) {
		// mock server
		mockedResponseBody := `{
		  "100.11.123.208": {
			"cc": "US",
			"country": "United States",
			"region": "North America",
			"city": "Lansdale",
			"coordinates": [ 40.2363, -75.296 ]
          }
	}`

		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, mockedResponseBody)
		}))
		defer mockServer.Close()

		storeFromUrl := geo.NewStoreFromUrl(mockServer.URL)
		assert.Equal(t, store, storeFromUrl)
	})
}
