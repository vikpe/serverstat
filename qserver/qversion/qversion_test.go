package qversion_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/qversion"
)

func TestVersion_IsMvdsv(t *testing.T) {
	testCases := map[string]bool{
		"mvdsv":     true,
		"MVDSV 1.2": true,
		"":          false,
		"foo":       false,
	}

	for version, expect := range testCases {
		assert.Equal(t, expect, qversion.New(version).IsMvdsv(), version)
	}
}

func TestVersion_IsFte(t *testing.T) {
	testCases := map[string]bool{
		"fte":     true,
		"FTE 1.2": true,
		"":        false,
		"foo":     false,
	}

	for version, expect := range testCases {
		assert.Equal(t, expect, qversion.New(version).IsFte(), version)
	}
}

func TestVersion_IsProxy(t *testing.T) {
	testCases := map[string]bool{
		"qwfwd":     true,
		"QWFWD 1.2": true,
		"":          false,
		"foo":       false,
	}

	for version, expect := range testCases {
		assert.Equal(t, expect, qversion.New(version).IsProxy(), version)
	}
}

func TestVersion_IsQtv(t *testing.T) {
	testCases := map[string]bool{
		"qtv":     true,
		"QTV 1.2": true,
		"":        false,
		"foo":     false,
	}

	for version, expect := range testCases {
		assert.Equal(t, expect, qversion.New(version).IsQtv(), version)
	}
}

func TestVersion_IsGameServer(t *testing.T) {
	testCases := map[string]bool{
		"fte":   true,
		"mvdsv": true,
		"qtv":   false,
		"qwfwd": false,
		"":      false,
		"foo":   false,
	}

	for version, expect := range testCases {
		assert.Equal(t, expect, qversion.New(version).IsGameServer(), version)
	}
}

func TestVersion_GetType(t *testing.T) {
	testCases := map[string]qversion.Type{
		"mvdsv":  qversion.TypeMvdsv,
		"qwfwd":  qversion.TypeProxy,
		"qtv":    qversion.TypeQtv,
		"fte":    qversion.TypeFte,
		"foobar": qversion.TypeUnknown,
	}

	for version, expect := range testCases {
		assert.Equal(t, expect, qversion.New(version).GetType(), version)
	}
}