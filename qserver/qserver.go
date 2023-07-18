package qserver

import (
	"fmt"
	"net"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/samber/lo"
	"github.com/vikpe/serverstat/qserver/geo"
	"github.com/vikpe/serverstat/qserver/mvdsv/qtvstream"
	"github.com/vikpe/serverstat/qserver/qclient"
	"github.com/vikpe/serverstat/qserver/qsettings"
	"github.com/vikpe/serverstat/qserver/qversion"
)

type GenericServer struct {
	Address   string             `json:"address"`
	Version   qversion.Version   `json:"version"`
	Clients   []qclient.Client   `json:"clients"`
	Settings  qsettings.Settings `json:"settings"`
	Geo       geo.Location       `json:"geo"`
	ExtraInfo struct {
		QtvStream qtvstream.QtvStream `json:"qtv_stream"`
	} `json:"extra_info"`
}

func (server GenericServer) Players() []qclient.Client {
	return lo.Filter(server.Clients, func(c qclient.Client, i int) bool {
		return c.IsPlayer()
	})
}

func (server GenericServer) Spectators() []qclient.Client {
	return lo.Filter(server.Clients, func(c qclient.Client, i int) bool {
		return c.IsSpectator()
	})
}

func ParseHostname(address string, hostname string) string {
	if "" == hostname || !strings.Contains(hostname, ".") {
		return address
	}

	potentialHostname := strings.ToLower(hostname)

	if strings.Contains(potentialHostname, " ") {
		hostnameParts := strings.Split(potentialHostname, " ")
		potentialHostname = strings.TrimSpace(hostnameParts[0])
	}

	if len(potentialHostname) < 4 || !strings.Contains(potentialHostname, ".") {
		return address
	}

	validate := validator.New()

	if strings.Contains(potentialHostname, ":") {
		commaIndex := strings.Index(potentialHostname, ":")
		const suffixLength = len(":28501")

		if len(potentialHostname) < (commaIndex + suffixLength) {
			return address
		}

		potentialHostname = potentialHostname[0 : commaIndex+suffixLength]
		err := validate.Var(potentialHostname, "hostname_port")
		if err == nil {
			return potentialHostname
		}
	} else {
		err := validate.Var(potentialHostname, "hostname")
		if err == nil {
			_, port, _ := net.SplitHostPort(address)
			return fmt.Sprintf("%s:%s", potentialHostname, port)
		}
	}

	return address
}
