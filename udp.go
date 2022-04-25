package serverstat

import (
	"bytes"
	"log"
	"net"
)

func UdpRequest(address string, statusPacket []byte, expectedHeader []byte) ([]byte, error) {
	nullResponse := make([]byte, 0)

	conn, err := net.Dial("udp4", address)
	if err != nil {
		return nullResponse, err
	}
	defer conn.Close()

	response := make([]byte, 8192)

	const (
		Retries     = 3
		TimeoutInMs = 500
	)

	for i := 0; i < Retries; i++ {
		conn.SetDeadline(timeInFuture(TimeoutInMs))

		_, err = conn.Write(statusPacket)
		if err != nil {
			return nullResponse, err
		}

		conn.SetDeadline(timeInFuture(TimeoutInMs))
		_, err = conn.Read(response)
		if err != nil {
			continue
		}

		break
	}

	if err != nil {
		return nullResponse, err
	}

	isValidHeader := bytes.Equal(response[:len(expectedHeader)], expectedHeader)
	if !isValidHeader {
		log.Println(address + ": Response error, invalid header.")
		return nullResponse, err
	}

	return response, nil
}
