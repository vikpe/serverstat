package qstatus

import (
	"strings"

	"golang.org/x/exp/slices"
)

const (
	Countdown = "Countdown"
	Started   = "Started"
	Standby   = "Standby"
	Unknown   = "Unknown"
)

func Parse(status string) string {
	knownStatuses := []string{Standby, Countdown, Started}

	if slices.Contains(knownStatuses, status) {
		return status
	} else if strings.Contains(status, " min left") {
		return Started
	} else {
		return Unknown
	}
}
