package mvdsv_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/mvdsv"
)

func TestQtvStream_MarshalJSON(t *testing.T) {
	// empty URL
	emptyStream := mvdsv.QtvStream{Url: ""}
	result, _ := json.Marshal(emptyStream)
	expect := `""`
	assert.Equal(t, expect, string(result))

	// non-empty URL
	stream := mvdsv.QtvStream{Url: "1@qw.foppa.dk:28000"}
	expect = `{"Title":"","Url":"1@qw.foppa.dk:28000","Clients":null,"NumClients":0}`
	result, _ = json.Marshal(stream)
	assert.Equal(t, expect, string(result))
}