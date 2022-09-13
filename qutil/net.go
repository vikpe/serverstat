package qutil

import (
	"net"

	"github.com/go-playground/validator/v10"
)

func HostnameToIp(hostname string) string {
	ips, err := net.LookupIP(hostname)

	if err == nil {
		for _, ip := range ips {
			if ipv4 := ip.To4(); ipv4 != nil {
				return ipv4.String()
			}
		}
	}

	return hostname
}

func IsValidServerAddress(address string) bool {
	validate := validator.New()
	err := validate.Var(address, "required,hostname_port")
	return err == nil
}
