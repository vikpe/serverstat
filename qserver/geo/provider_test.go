package geo_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/geo"
)

func TestMemcachedProvider(t *testing.T) {
	// mock server
	requestCount := 0
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
		requestCount++
	}))
	defer mockServer.Close()

	cacheDuration := time.Millisecond * 50
	provider := geo.NewMemcachedProvider(mockServer.URL, cacheDuration)
	expectedLocation := geo.Location{
		Region:      "North America",
		Country:     "United States",
		CC:          "US",
		City:        "Lansdale",
		Coordinates: [2]float32{40.2363, -75.296},
	}

	// cache
	// initial request
	assert.Equal(t, expectedLocation, provider.ByIp("100.11.123.208"))
	assert.Equal(t, 1, requestCount)

	// no subsequent request (use cache)
	assert.Equal(t, expectedLocation, provider.ByIp("100.11.123.208"))
	assert.Equal(t, 1, requestCount)

	// new request (expired cache)
	time.Sleep(cacheDuration + time.Millisecond)
	assert.Equal(t, expectedLocation, provider.ByIp("100.11.123.208"))
	assert.Equal(t, 2, requestCount)

	// all methods
	assert.Equal(t, geo.Location{}, provider.ByIp("foo"))
	assert.Equal(t, expectedLocation, provider.ByAddress("100.11.123.208:28501"))
	assert.Equal(t, expectedLocation, provider.ByHostname("100.11.123.208"))
}
