package eos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitializeAndClose(t *testing.T) {
	e := NewEOSClient()
	assert.Nil(t, e.Initialize())
	assert.Nil(t, e.Close())
}
