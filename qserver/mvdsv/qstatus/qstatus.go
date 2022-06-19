package qstatus

import (
	"fmt"
	"strings"

	"github.com/vikpe/serverstat/qserver/mvdsv/qmode"
	"github.com/vikpe/serverstat/qserver/qclient"
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

func New(status string, mode qmode.Mode, players []qclient.Client, freeSlots int) Status {
	if mode.IsRace() {
		return Status{
			Name:        Standby,
			Description: "Racing",
		}
	}

	if Standby == status && (mode.IsXonX() || mode.IsFfa()) && hasFrags(players) {
		if hasBots(players) {
			if hasHumans(players) {
				return Status{Name: Standby, Description: "Score screen"}
			}

			return Status{Name: Standby, Description: "Waiting for players"}
		}

		return Status{Name: Started, Description: "Score screen"}
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

func hasFrags(players []qclient.Client) bool {
	if 0 == len(players) {
		return false
	}

	for _, p := range players {
		if p.Frags > 0 {
			return true
		}
	}

	return false
}

func hasBots(players []qclient.Client) bool {
	if 0 == len(players) {
		return false
	}

	for _, p := range players {
		if p.IsBot() {
			return true
		}
	}

	return false
}

func hasHumans(players []qclient.Client) bool {
	if 0 == len(players) {
		return false
	}

	for _, p := range players {
		if p.IsHuman() {
			return true
		}
	}

	return false
}
