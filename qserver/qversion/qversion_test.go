package qversion_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/serverstat/qserver/qversion"
)

func TestVersion_IsMvdsv(t *testing.T) {
	assert.True(t, qversion.New("mvdsv").IsMvdsv())
	assert.True(t, qversion.New("MVDSV 0.35-dev").IsMvdsv())
	assert.False(t, qversion.New("").IsMvdsv())
	assert.False(t, qversion.New("foo").IsMvdsv())
}

func TestVersion_IsFte(t *testing.T) {
	assert.True(t, qversion.New("fte").IsFte())
	assert.True(t, qversion.New("fte 1.2").IsFte())
	assert.False(t, qversion.New("").IsFte())
	assert.False(t, qversion.New("foo").IsFte())
}

func TestVersion_IsProxy(t *testing.T) {
	assert.True(t, qversion.New("qwfwd").IsProxy())
	assert.True(t, qversion.New("qwfwd 1.2").IsProxy())
	assert.False(t, qversion.New("").IsProxy())
	assert.False(t, qversion.New("foo").IsProxy())
}

func TestVersion_IsQtv(t *testing.T) {
	assert.True(t, qversion.New("qtv").IsQtv())
	assert.True(t, qversion.New("qtv 1.2").IsQtv())
	assert.False(t, qversion.New("").IsQtv())
	assert.False(t, qversion.New("foo").IsQtv())
}

func TestVersion_IsGameServer(t *testing.T) {
	assert.True(t, qversion.New("fte").IsGameServer())
	assert.True(t, qversion.New("fte 1.2").IsGameServer())
	assert.True(t, qversion.New("mvdsv").IsGameServer())
	assert.True(t, qversion.New("MVDSV 0.35-dev").IsGameServer())

	assert.False(t, qversion.New("").IsGameServer())
	assert.False(t, qversion.New("foo").IsGameServer())
}

func TestVersion_GetType(t *testing.T) {
	assert.Equal(t, qversion.TypeMvdsv, qversion.New("mvdsv 1.0").GetType())
	assert.Equal(t, qversion.TypeProxy, qversion.New("qwfwd 1.0").GetType())
	assert.Equal(t, qversion.TypeQtv, qversion.New("qtv 1.0").GetType())
	assert.Equal(t, qversion.TypeFte, qversion.New("fte 1.0").GetType())
	assert.Equal(t, qversion.TypeUnknown, qversion.New("foobar").GetType())
}
