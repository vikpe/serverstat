package qwnet

import (
	"bytes"
	"errors"
	"net"
	"time"
)

const (
	defaultBufferSize  uint16 = 8192
	defaultRetries     uint8  = 3
	defaultTimeoutInMs uint16 = 500
)

type UdpClientConfig struct {
	BufferSize  uint16
	Retries     uint8
	TimeoutInMs uint16
}

type UdpClient struct {
	Config UdpClientConfig
}

func NewUdpClient() *UdpClient {
	defaultConfig := UdpClientConfig{
		BufferSize:  defaultBufferSize,
		Retries:     defaultRetries,
		TimeoutInMs: defaultTimeoutInMs,
	}
	return NewUdpClientWithConfig(defaultConfig)
}

func NewUdpClientWithConfig(config UdpClientConfig) *UdpClient {
	return &UdpClient{
		Config: config,
	}
}

func (client UdpClient) getTimeout() time.Time {
	return time.Now().Add(time.Duration(client.Config.TimeoutInMs) * time.Millisecond)
}

func (client UdpClient) Request(address string, statusPacket []byte, expectedResponseHeader []byte) ([]byte, error) {
	conn, err := net.Dial("udp4", address)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	response := make([]byte, client.Config.BufferSize)
	responseLength := 0

	for i := uint8(0); i < client.Config.Retries; i++ {
		conn.SetDeadline(client.getTimeout())

		_, err = conn.Write(statusPacket)
		if err != nil {
			return nil, err
		}

		conn.SetDeadline(client.getTimeout())
		responseLength, err = conn.Read(response)
		if err != nil {
			continue
		}

		break
	}

	if err != nil {
		return nil, err
	}

	isValidResponseHeader := bytes.Equal(response[:len(expectedResponseHeader)], expectedResponseHeader)
	if !isValidResponseHeader {
		err = errors.New(address + ": Invalid response header.")
		return nil, err
	}

	return response[:responseLength], nil
}
