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
	if mode.IsRace() {
		return Status{
			Name:        Standby,
			Description: "Racing",
		}
	}

	if Standby == status {
		var description string

		if freeSlots > 0 && mode.IsXonX() {
			description = fmt.Sprintf("Waiting for %d %s", freeSlots, qutil.Pluralize("player", freeSlots))
		} else {
			description = "Waiting for players to ready up"
		}

		return Status{
			Name:        status,
			Description: description,
		}
	} else if Countdown == status || strings.Contains(status, " min left") {
		var description string

		if mode.IsCoop() {
			description = "Game in progress"
		} else {
			description = status
		}

		return Status{
			Name:        Started,
			Description: description,
		}
	}

	return Status{
		Name:        Unknown,
		Description: status,
	}
}
