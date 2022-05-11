package qstatus

import (
	"strings"

	"github.com/vikpe/serverstat/qserver/qsettings"
	"github.com/vikpe/serverstat/qutil"
)

const (
	statusStarted = "Started"
	stausStandBy  = "Standby"
)

type MatchDuration struct {
	Elapsed   int
	Total     int
	Remaining int
}

type Status struct {
	Name     string
	Duration MatchDuration
}

func Parse(settings qsettings.Settings) Status {
	status := settings.Get("status", "")
	const minLeftNeedle = " min left"

	var name string
	timeTotal := settings.GetInt("timelimit", 0)
	timeRemaining := timeTotal

	if strings.Contains(status, minLeftNeedle) {
		name = statusStarted
		timeRemaining = qutil.StringToInt(strings.Replace(status, minLeftNeedle, "", 1))
	} else {
		name = stausStandBy
	}

	return Status{
		Name: name,
		Duration: MatchDuration{
			Total:     timeTotal,
			Elapsed:   timeTotal - timeRemaining,
			Remaining: timeRemaining,
		},
	}
}
