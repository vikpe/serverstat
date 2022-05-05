package qsettings

import "strings"

func New(settingsString string) map[string]string {
	settingsLines := strings.FieldsFunc(settingsString, func(r rune) bool { return r == '\\' })
	settings := make(map[string]string, 0)

	for i := 0; i < len(settingsLines)-1; i += 2 {
		settings[settingsLines[i]] = settingsLines[i+1]
	}

	return settings
}
