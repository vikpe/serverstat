package qwnet_test

import "github.com/vikpe/qw-serverstat/qwnet"

func ExampleUdpRequest() {
	statusPacket := []byte{0xff, 0xff, 0xff, 0xff, 's', 't', 'a', 't', 'u', 's', ' ', '2', '3', 0x0a}
	expectedHeader := []byte{0xff, 0xff, 0xff, 0xff, 'n', '\\'}
	response, err := qwnet.UdpRequest("qw.foppa.dk:27502", statusPacket, expectedHeader)
}
