package laststats

// todo
/*import (
	"fmt"

	"github.com/vikpe/serverstat/qtext/qstring"
	"github.com/vikpe/udpclient"
)

var Command = udpclient.Command{
	RequestPacket:  []byte{0xff, 0xff, 0xff, 0xff, 'l', 'a', 's', 't', 's', 'c', 'o', 'r', 'e', 's', ' ', '2', 0x0a},
	ResponseHeader: []byte{0xff, 0xff, 0xff, 0xff, 'n', 'L', 'i', 's', 't', ' ', 'o', 'f', ' ', '2', ' ', 'l', 'a', 's', 't', ' ', 'd', 'e', 'm', 'o', 's', ':'},
}

func ParseResponse(responseBody []byte, err error) ([]Entry, error) {
	if err != nil {
		fmt.Print("ParseResponse err", string(responseBody))
		return make([]Entry, 0), err
	}

	fmt.Println(qstring.New(string(responseBody)).ToPlainString())

	/*rows := strings.Split(string(responseBody), "\n")
	//fmt.Println(rows)

	for _, r := range rows {
		fmt.Println(qstring.New(r).ToPlainString())
	}

	spectatorPlainNames := make([]Entry, 0)

	return spectatorPlainNames, nil
}*/
