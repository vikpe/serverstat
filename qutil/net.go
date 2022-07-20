package qutil

import "net"

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
