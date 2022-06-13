package qstatus

import (
	"fmt"
	"strings"

	"github.com/vikpe/serverstat/qserver/mvdsv/qmode"
	"github.com/vikpe/serverstat/qutil"
)

const (
	Countdown = "Countdown"
	Started   = "Started"
	Standby   = "Standby"
	Unknown   = "Unknown"
)

type Status struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func New(status string, freeSlots int, mode qmode.Mode) Status {
	if Standby == status {
		description := "Waiting for players to ready up"

		if 0 != freeSlots && mode.IsXonX() {
			description = fmt.Sprintf("Waiting for %d %s", freeSlots, qutil.Pluralize("player", freeSlots))
		}

		return Status{
			Name:        status,
			Description: description,
		}
	} else if Countdown == status || strings.Contains(status, " min left") {
		return Status{
			Name:        Started,
			Description: status,
		}
	}

	return Status{
		Name:        Unknown,
		Description: status,
	}
}
