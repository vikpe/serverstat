package qtime

import (
	"strings"

	"github.com/vikpe/serverstat/qutil"
)

type Time struct {
	Elapsed   int `json:"elapsed"`
	Total     int `json:"total"`
	Remaining int `json:"remaining"`
}

func Parse(timelimit int, status string) Time {
	timeRemaining := timelimit

	const minLeftNeedle = " min left"

	if strings.Contains(status, minLeftNeedle) {
		timeRemaining = qutil.StringToInt(strings.Replace(status, minLeftNeedle, "", 1))
	}

	return Time{
		Total:     timelimit,
		Elapsed:   timelimit - timeRemaining,
		Remaining: timeRemaining,
	}
}
