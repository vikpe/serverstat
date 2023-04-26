package laststats_test

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/mvdsv/commands/laststats"
)

func TestGetCommand(t *testing.T) {
	command := laststats.GetCommand(5)
	assert.Equal(t, "\xff\xff\xff\xfflaststats 5\n", string(command.RequestPacket))
	assert.Equal(t, "\xff\xff\xff\xffnlaststats ", string(command.ResponseHeader))
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
				assert.ErrorContains(t, err, "malformed json")
			})

			t.Run("not containing json", func(t *testing.T) {
				stats, err := laststats.ParseResponseBody([]byte("����nlaststats 2\nfoo"), nil)
				assert.Equal(t, []laststats.Entry{}, stats)
				assert.ErrorContains(t, err, "malformed json")
			})

			t.Run("invalid json", func(t *testing.T) {
				t.Run("invalid format 1", func(t *testing.T) {
					stats, err := laststats.ParseResponseBody([]byte("����nlaststats 2\n["), nil)
					assert.Equal(t, []laststats.Entry{}, stats)
					assert.ErrorContains(t, err, "malformed json")
				})

				t.Run("invalid format 2", func(t *testing.T) {
					stats, err := laststats.ParseResponseBody([]byte("����nlaststats 2\n]["), nil)
					assert.Equal(t, []laststats.Entry{}, stats)
					assert.ErrorContains(t, err, "malformed json")
				})

				t.Run("invalid format 3", func(t *testing.T) {
					stats, err := laststats.ParseResponseBody([]byte("����nlaststats 2\n]"), nil)
					assert.Equal(t, []laststats.Entry{}, stats)
					assert.ErrorContains(t, err, "malformed json")
				})

				t.Run("invalid format 4", func(t *testing.T) {
					stats, err := laststats.ParseResponseBody([]byte(`����nlaststats 2\n[{]`), nil)
					assert.Equal(t, []laststats.Entry{}, stats)
					assert.ErrorContains(t, err, "invalid character")
				})

				t.Run("invalid fields/values in json", func(t *testing.T) {
					stats, err := laststats.ParseResponseBody([]byte(`����nlaststats 2\n[{"foo": "bar"}]`), nil)
					assert.Equal(t, []laststats.Entry{}, stats)
					assert.ErrorContains(t, err, "invalid fields, date is missing")
				})
			})
		})

		t.Run("valid response body", func(t *testing.T) {
			t.Run("empty stats", func(t *testing.T) {
				responseBody := []byte("����nlaststats 0\n[]\n")
				stats, err := laststats.ParseResponseBody(responseBody, nil)
				assert.Equal(t, []laststats.Entry{}, stats)
				assert.Nil(t, err)
			})

			t.Run("non-empty stats", func(t *testing.T) {
				responseBody, _ := os.ReadFile("./test_files/response.bin")
				stats, err := laststats.ParseResponseBody(responseBody, nil)
				assert.Len(t, stats, 2, err)
				assert.Equal(t, "2023-04-25 21:12:16 +0200", stats[0].Date)
				assert.Equal(t, "2023-04-25 21:22:49 +0200", stats[1].Date)
				assert.Nil(t, err)
			})
		})
	})
}
