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

// A single Canon T4i camera must be connected in order to run as expected.
func TestGetCameraModels(t *testing.T) {
	e := NewEOSClient()
	e.Initialize()
	defer e.Close()

	models, err := e.GetCameraModels()
	assert.NotNil(t, models)
	assert.Equal(t, 1, len(models))
	assert.Nil(t, err)

	// verify values in the first model entry
	assert.Equal(t, "Canon EOS REBEL T4i", models[0].szDeviceDescription)
	assert.Equal(t, "0", models[0].szPortName)
	assert.Equal(t, 2971958586, int(models[0].reserved))
	assert.Equal(t, 1, int(models[0].deviceSubType))
}
