package eos

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework EDSDK
#define __MACOS__ 1
#include <EDSDK/EDSDK.h>
#include <EDSDK/EDSDKTypes.h>
#include <stdlib.h>
*/
import (
	"C"
)
import (
	"errors"
	"fmt"
	"unsafe"
)

type LiveViewOutputDevice int

const (
	TFT LiveViewOutputDevice = iota
	PC
)

type CameraModel struct {
	camera         *C.EdsCameraRef
	sessionOpen    bool
	liveViewActive bool
	liveViewDevice int

	szPortName          string
	szDeviceDescription string
	deviceSubType       uint32
	reserved            uint32
}

// Releases reference to the camera
func (c *CameraModel) Release() {
	C.EdsRelease((*C.struct___EdsObject)(unsafe.Pointer(&c.camera)))
}

// Open a session with the camera for sending commands
func (c *CameraModel) OpenSession() error {
	eosError := C.EdsOpenSession((*C.struct___EdsObject)(unsafe.Pointer(c.camera)))
	if eosError != C.EDS_ERR_OK {
		return errors.New(fmt.Sprintf("Error when opening session with camera (code=%d)", eosError))
	}
	c.sessionOpen = true
	return nil
}

// Close an existing camera session
func (c *CameraModel) CloseSession() {
	if c.sessionOpen == false {
		return
	}
	c.sessionOpen = false
	C.EdsCloseSession((*C.struct___EdsObject)(unsafe.Pointer(&c.camera)))
}

// Take a picture
func (c *CameraModel) TakePicture() error {
	if c.sessionOpen == false {
		return errors.New("Session is not open, must call OpenSession first")
	}
	eosError := C.EdsSendCommand((*C.struct___EdsObject)(unsafe.Pointer(c.camera)), C.kEdsCameraCommand_TakePicture, 0)
	if eosError != C.EDS_ERR_OK {
		return errors.New(fmt.Sprintf("Error when taking picture (code=%d)", eosError))
	}
	return nil
}

// Start LiveView on the device configured with SetLiveViewOutputDevice
func (c *CameraModel) StartLiveView() error {
	if c.sessionOpen == false {
		return errors.New("Session is not open, must call OpenSession first")
	}

	if c.liveViewActive == true {
		return errors.New("LiveView is already active, cannot start")
	}

	var device int
	eosError := C.EdsGetPropertyData((*C.struct___EdsObject)(unsafe.Pointer(c.camera)), C.kEdsPropID_Evf_OutputDevice, 0, (C.EdsUInt32)(unsafe.Sizeof(device)), unsafe.Pointer(&device))
	if eosError != C.EDS_ERR_OK {
		return errors.New(fmt.Sprintf("Error getting output device property when activating LiveMode (code=%d)", eosError))
	}

	// connect Live View output device
	device |= c.liveViewDevice
	eosError = C.EdsSetPropertyData((*C.struct___EdsObject)(unsafe.Pointer(c.camera)), C.kEdsPropID_Evf_OutputDevice, 0, (C.EdsUInt32)(unsafe.Sizeof(device)), unsafe.Pointer(&device))
	if eosError != C.EDS_ERR_OK {
		return errors.New(fmt.Sprintf("Error setting output device property when activating LiveMode (code=%d)", eosError))
	}
	c.liveViewActive = true
	return nil
}

// Stop LiveView on the device configured with SetLiveViewOutputDevice
func (c *CameraModel) StopLiveView() error {
	if c.sessionOpen == false {
		return errors.New("Session is not open, must call OpenSession first")
	}

	if c.liveViewActive == false {
		return errors.New("LiveView is already inactive, cannot stop")
	}

	var device int
	eosError := C.EdsGetPropertyData((*C.struct___EdsObject)(unsafe.Pointer(c.camera)), C.kEdsPropID_Evf_OutputDevice, 0, (C.EdsUInt32)(unsafe.Sizeof(device)), unsafe.Pointer(&device))
	if eosError != C.EDS_ERR_OK {
		return errors.New(fmt.Sprintf("Error getting output device property when stopping LiveMode (code=%d)", eosError))
	}

	// disconnect Live View output device
	device &= ^c.liveViewDevice
	eosError = C.EdsSetPropertyData((*C.struct___EdsObject)(unsafe.Pointer(c.camera)), C.kEdsPropID_Evf_OutputDevice, 0, (C.EdsUInt32)(unsafe.Sizeof(device)), unsafe.Pointer(&device))
	if eosError != C.EDS_ERR_OK {
		return errors.New(fmt.Sprintf("Error setting output device property when stopping LiveMode (code=%d)", eosError))
	}
	c.liveViewActive = false
	return nil
}

// Toggle the LiveView state of the camera
func (c *CameraModel) ToggleLiveView() error {
	if c.liveViewActive {
		return c.StopLiveView()
	} else {
		return c.StartLiveView()
	}
}

// Set the device to use with LiveView.  Will stop LiveView if already active on a device
func (c *CameraModel) SetLiveViewOutputDevice(device LiveViewOutputDevice) error {
	if c.liveViewActive {
		if err := c.StopLiveView(); err != nil {
			return err
		}
	}

	switch device {
	case TFT:
		c.liveViewDevice = C.kEdsEvfOutputDevice_TFT
		break
	case PC:
		c.liveViewDevice = C.kEdsEvfOutputDevice_PC
		break
	default:
		return errors.New("Unrecognized LiveView device supplied")
	}
	return nil
}
