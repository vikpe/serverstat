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

func TestVersion_IsQwfwd(t *testing.T) {
	testCases := map[string]bool{
		"qwfwd":     true,
		"QWFWD 1.2": true,
		"":          false,
		"foo":       false,
	}

	for version, expect := range testCases {
		assert.Equal(t, expect, qversion.New(version).IsQwfwd(), version)
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

func TestVersion_IsFortressOne(t *testing.T) {
	testCases := map[string]bool{
		"fo svn":      true,
		"FO SVN 6128": true,
		"":            false,
		"foo":         false,
	}

	for version, expect := range testCases {
		assert.Equal(t, expect, qversion.New(version).IsFortressOne(), version)
	}
}

func TestVersion_GetType(t *testing.T) {
	testCases := map[string]string{
		"Mvdsv":  "mvdsv",
		"QWFWD":  "qwfwd",
		"qtv":    "qtv",
		"fte":    "fte",
		"fo svn": "fortress_one",
		"foobar": "unknown",
	}

	for version, expectType := range testCases {
		assert.Equal(t, expectType, qversion.New(version).GetType(), version)
	}
}
