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

func TestSetCameraAddedHandler(t *testing.T) {
	e := NewEOSClient()
	e.Initialize()
	defer e.Close()

	f := func() { t.Log("Camera connected!") }
	e.SetCameraAddedHandler(f)
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
	camera := models[0]
	defer camera.Release()
	assert.NotNil(t, camera.camera)
	assert.Equal(t, "Canon EOS REBEL T4i", camera.szDeviceDescription)
	assert.Equal(t, "0", camera.szPortName)
	assert.Equal(t, 2971958586, int(camera.reserved))
	assert.Equal(t, 1, int(camera.deviceSubType))
}
