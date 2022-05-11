package qsettings

import (
	"strconv"
	"strings"
)

type Settings map[string]string

func (settings Settings) Has(key string) bool {
	_, hasKey := settings[key]
	return hasKey
}

func (settings Settings) Get(key string, default_ string) string {
	if value, ok := settings[key]; ok {
		return value
	} else {
		return default_
	}
}

func (settings Settings) GetInt(key string, default_ int) int {
	if stringVal, ok := settings[key]; ok {
		intVal, _ := strconv.Atoi(stringVal)
		return intVal
	} else {
		return default_
	}
}

func ParseString(settingsString string) Settings {
	settingsLines := strings.FieldsFunc(settingsString, func(r rune) bool { return r == '\\' })
	settings := Settings{}

	for i := 0; i < len(settingsLines)-1; i += 2 {
		settings[settingsLines[i]] = settingsLines[i+1]
	}

	return settings
}
