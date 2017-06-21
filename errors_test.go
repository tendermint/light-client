package lightclient

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorHeight(t *testing.T) {
	e1 := ErrHeightMismatch(2, 3)
	e1.Error()
	assert.True(t, IsHeightMismatchErr(e1))

	e2 := errors.New("foobar")
	assert.False(t, IsHeightMismatchErr(e2))
	assert.False(t, IsHeightMismatchErr(nil))
}

func TestErrorNoData(t *testing.T) {
	e1 := ErrNoData()
	e1.Error()
	assert.True(t, IsNoDataErr(e1))

	e2 := errors.New("foobar")
	assert.False(t, IsNoDataErr(e2))
	assert.False(t, IsNoDataErr(nil))
}
