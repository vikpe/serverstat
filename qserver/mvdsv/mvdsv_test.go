package mvdsv_test

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/mvdsv/lastscores"
	"github.com/vikpe/serverstat/qserver/mvdsv/qtvstream"
	"github.com/vikpe/udphelper"
)

func TestGetQtvStream(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		stream, err := mvdsv.GetQtvStream("foo:666")
		assert.Equal(t, qtvstream.QtvStream{
			SpectatorNames: make([]string, 0),
		}, stream)
		assert.NotNil(t, err)
	})

	t.Run("has no clients", func(t *testing.T) {
		go func() {
			qtvHeader := []byte{0xff, 0xff, 0xff, 0xff, 'n', 'q', 't', 'v', ' '}
			qtvBody := []byte(`1 "qw.foppa.dk - qtv (3)" "3@qw.foppa.dk:28000" 0`)
			qtvResponse := append(qtvHeader, qtvBody...)

			udphelper.New(":5001").Respond(qtvResponse)
		}()
		time.Sleep(10 * time.Millisecond)

		stream, err := mvdsv.GetQtvStream(":5001")
		expectStream := qtvstream.QtvStream{
			Title:          "qw.foppa.dk - qtv (3)",
			Url:            "3@qw.foppa.dk:28000",
			ID:             3,
			Address:        "qw.foppa.dk:28000",
			SpectatorCount: 0,
			SpectatorNames: []string{},
		}
		assert.Equal(t, expectStream, stream)
		assert.Nil(t, err)
	})

	t.Run("has clients", func(t *testing.T) {
		go func() {
			qtvHeader := []byte{0xff, 0xff, 0xff, 0xff, 'n', 'q', 't', 'v', ' '}
			qtvBody := []byte(`1 "qw.foppa.dk - qtv (3)" "3@qw.foppa.dk:28000" 2`)
			qtvResponse := append(qtvHeader, qtvBody...)
			qtvusersHeader := []byte{0xff, 0xff, 0xff, 0xff, 'n', 'q', 't', 'v', 'u', 's', 'e', 'r', 's'}
			bassInRed := string([]byte{98 + 128, 97 + 128, 115 + 128, 115 + 128})
			qtvusersBody := []byte(fmt.Sprintf(`12 "XantoM" "%s"`, bassInRed))
			qtvusersResponse := append(qtvusersHeader, qtvusersBody...)

			udphelper.New(":5002").Respond(qtvResponse, qtvusersResponse)
		}()
		time.Sleep(10 * time.Millisecond)

		stream, err := mvdsv.GetQtvStream(":5002")
		expectStream := qtvstream.QtvStream{
			Title:          "qw.foppa.dk - qtv (3)",
			Url:            "3@qw.foppa.dk:28000",
			ID:             3,
			Address:        "qw.foppa.dk:28000",
			SpectatorCount: 2,
			SpectatorNames: []string{
				"XantoM",
				"bass",
			},
		}

		assert.Equal(t, expectStream, stream)
		assert.Nil(t, err)
	})
}

func TestGetLastScores(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		scores, err := mvdsv.GetLastScores("foo:666", 5)
		assert.Equal(t, []lastscores.Entry{}, scores)
		assert.NotNil(t, err)
	})

	t.Run("success", func(t *testing.T) {
		t.Run("empty result", func(t *testing.T) {
			go func() {
				response := []byte{0xff, 0xff, 0xff, 0xff, 'n', 'l', 'a', 's', 't', 's', 't', 'a', 't', 's', ' ', '2', 0xa, '[', ']'}
				udphelper.New(":5003").Respond(response)
			}()
			time.Sleep(10 * time.Millisecond)

			scores, err := mvdsv.GetLastScores(":5003", 5)
			assert.Equal(t, scores, []lastscores.Entry{})
			assert.Nil(t, err)
		})

		t.Run("non-empty result", func(t *testing.T) {
			go func() {
				responseHeader := []byte{0xff, 0xff, 0xff, 0xff, 'n', 'l', 'a', 's', 't', 's', 't', 'a', 't', 's', ' ', '2', 0xa}
				responseBody, _ := ioutil.ReadFile("./commands/laststats/test_files/laststats.json")
				response := append(responseHeader, responseBody...)
				udphelper.New(":5004").Respond(response)
			}()
			time.Sleep(10 * time.Millisecond)

			scores, err := mvdsv.GetLastScores(":5004", 5)

			fmt.Println("err", err)

			assert.Len(t, scores, 2)
			assert.Equal(t, "duel", scores[0].Mode)
			assert.Equal(t, "4on4", scores[1].Mode)
			assert.Nil(t, err)
		})
	})
}
