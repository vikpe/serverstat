package geo_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/geo"
)

func TestNewStoreFromUrl(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		store := geo.NewStoreFromUrl("foo")
		assert.Equal(t, geo.Store{}, store)
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

		// test methods
		store := geo.NewStoreFromUrl(mockServer.URL)

		expectedLocation := geo.Location{
			CC:          "US",
			Country:     "United States",
			Region:      "North America",
			City:        "Lansdale",
			Coordinates: [2]float32{40.2363, -75.296},
		}

		t.Run("ByIp", func(t *testing.T) {
			t.Run("known", func(t *testing.T) {
				assert.Equal(t, expectedLocation, store.ByIp("100.11.123.208"))
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
			assert.Equal(t, expectedLocation, store.ByAddress("100.11.123.208:28000"))
		})
	})

}
