package laststats_test

import (
	"errors"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/mvdsv/commands/laststats"
	"github.com/vikpe/udpclient"
)

// todo
func TestGetCommand(t *testing.T) {
	expect := udpclient.Command{
		RequestPacket:  []byte{0xff, 0xff, 0xff, 0xff, 'l', 'a', 's', 't', 's', 't', 'a', 't', 's', ' ', byte(5), 0x0a},
		ResponseHeader: []byte{0xff, 0xff, 0xff, 0xff, 'n', 'l', 'a', 's', 't', 's', 't', 'a', 't', 's', ' '},
	}
	assert.Equal(t, expect, laststats.GetCommand(5))
}

func TestParseResponseBody(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		result, err := laststats.ParseResponseBody([]byte{}, errors.New("some error"))
		assert.Equal(t, []laststats.Entry{}, result)
		assert.ErrorContains(t, err, "some error")
	})

	t.Run("no error", func(t *testing.T) {
		t.Run("invalid response body", func(t *testing.T) {
			t.Run("empty", func(t *testing.T) {
				stats, err := laststats.ParseResponseBody([]byte(""), nil)
				assert.Equal(t, []laststats.Entry{}, stats)
				assert.ErrorContains(t, err, "invalid json in response body")
			})

			t.Run("not containing json", func(t *testing.T) {
				stats, err := laststats.ParseResponseBody([]byte("2\nfoo"), nil)
				assert.Equal(t, []laststats.Entry{}, stats)
				assert.ErrorContains(t, err, "invalid json in response body")
			})

			t.Run("invalid json", func(t *testing.T) {
				t.Run("invalid format 1", func(t *testing.T) {
					stats, err := laststats.ParseResponseBody([]byte("2\n["), nil)
					assert.Equal(t, []laststats.Entry{}, stats)
					assert.ErrorContains(t, err, "invalid json in response body")
				})

				t.Run("invalid format 2", func(t *testing.T) {
					stats, err := laststats.ParseResponseBody([]byte("2\n]["), nil)
					assert.Equal(t, []laststats.Entry{}, stats)
					assert.ErrorContains(t, err, "invalid json in response body")
				})

				t.Run("invalid format 3", func(t *testing.T) {
					stats, err := laststats.ParseResponseBody([]byte("2\n]"), nil)
					assert.Equal(t, []laststats.Entry{}, stats)
					assert.ErrorContains(t, err, "invalid json in response body")
				})

				t.Run("invalid format 4", func(t *testing.T) {
					stats, err := laststats.ParseResponseBody([]byte(`2\n[{]`), nil)
					assert.Equal(t, []laststats.Entry{}, stats)
					assert.ErrorContains(t, err, "invalid json in response body")
				})

				t.Run("invalid fields/values in json", func(t *testing.T) {
					stats, err := laststats.ParseResponseBody([]byte(`2\n[{"foo": "bar"}]`), nil)
					assert.Equal(t, []laststats.Entry{}, stats)
					assert.ErrorContains(t, err, "invalid json in response body")
				})
			})
		})

		t.Run("valid response body", func(t *testing.T) {
			t.Run("empty stats", func(t *testing.T) {
				responseBody := []byte("0\n[]")
				stats, err := laststats.ParseResponseBody(responseBody, nil)
				assert.Equal(t, []laststats.Entry{}, stats)
				assert.Nil(t, err)
			})

			t.Run("non-empty stats", func(t *testing.T) {
				responseBodyHeader := []byte("2\n")
				responseBodyContent, _ := ioutil.ReadFile("test_files/laststats.json")
				responseBody := append(responseBodyHeader, responseBodyContent...)
				stats, err := laststats.ParseResponseBody(responseBody, nil)

				assert.Len(t, stats, 2)
				assert.Equal(t, "2022-10-31 19:59:46 +0100", stats[0].Date)
				assert.Equal(t, "2022-10-31 20:01:46 +0100", stats[1].Date)
				assert.Nil(t, err)
			})
		})
	})
}
