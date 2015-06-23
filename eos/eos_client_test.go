package eos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitializeAndClose(t *testing.T) {
	e := NewEOSClient()
	assert.Nil(t, e.Initialize())
	assert.Nil(t, e.Release())
}

// A single Canon T4i camera must be connected in order to run as expected.
func TestGetCameraModels(t *testing.T) {
	e := NewEOSClient()
	e.Initialize()
	defer e.Release()

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

// At least one camera must be connected in order to run successfully.
func TestTakePicture(t *testing.T) {
	e := NewEOSClient()
	e.Initialize()
	defer e.Release()

	models, _ := e.GetCameraModels()
	camera := models[0]
	defer camera.Release()
	assert.Nil(t, camera.OpenSession())
	defer camera.CloseSession()
	assert.Nil(t, camera.TakePicture())
}
