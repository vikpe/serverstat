package mvdsv_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/mvdsv"
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
