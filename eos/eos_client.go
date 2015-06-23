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

type EOSClient struct{}

func NewEOSClient() *EOSClient {
	return &EOSClient{}
}

// Initialize a new EOSClient instance
func (e *EOSClient) Initialize() error {
	eosError := C.EdsInitializeSDK()
	if eosError != C.EDS_ERR_OK {
		return errors.New("Error when initializing Canon SDK")
	}
	return nil
}

// Release the EOSClient, must be called on termination
func (p *EOSClient) Release() error {
	eosError := C.EdsTerminateSDK()
	if eosError != C.EDS_ERR_OK {
		return errors.New("Error when terminating Canon SDK")
	}
	return nil
}

// Get an array representing the cameras currently connected.  Each CameraModel
// instance must be released by invoking the Release function once no longer
// needed.
func (e *EOSClient) GetCameraModels() ([]CameraModel, error) {
	var eosCameraList *C.EdsCameraListRef
	var eosError C.EdsError
	var cameraCount int

	// get a reference to the cameras list record
	if eosError = C.EdsGetCameraList((*C.EdsCameraListRef)(unsafe.Pointer(&eosCameraList))); eosError != C.EDS_ERR_OK {
		return nil, errors.New(fmt.Sprintf("Error when obtaining reference to list of cameras (code=%d)", eosError))
	}
	defer C.EdsRelease((*C.struct___EdsObject)(unsafe.Pointer(&eosCameraList)))

	// get the number of cameras connected
	if eosError = C.EdsGetChildCount((*C.struct___EdsObject)(unsafe.Pointer(eosCameraList)), (*C.EdsUInt32)(unsafe.Pointer(&cameraCount))); eosError != C.EDS_ERR_OK {
		return nil, errors.New(fmt.Sprintf("Error when obtaining count of connected cameras (code=%d)", eosError))
	}

	// get details for each camera detected
	cameras := make([]CameraModel, 0)
	for i := 0; i < cameraCount; i++ {
		var eosCameraRef *C.EdsCameraRef
		if eosError = C.EdsGetChildAtIndex((*C.struct___EdsObject)(unsafe.Pointer(eosCameraList)), (C.EdsInt32)(i), (*C.EdsBaseRef)(unsafe.Pointer(&eosCameraRef))); eosError != C.EDS_ERR_OK {
			return nil, errors.New(fmt.Sprintf("Error when obtaining reference to camera (code=%d)", eosError))
		}

		var eosCameraDeviceInfo C.EdsDeviceInfo
		if eosError = C.EdsGetDeviceInfo((*C.struct___EdsObject)(unsafe.Pointer(eosCameraRef)), (*C.EdsDeviceInfo)(unsafe.Pointer(&eosCameraDeviceInfo))); eosError != C.EDS_ERR_OK {
			return nil, errors.New(fmt.Sprintf("Error when obtaining camera device info (code=%d)", eosError))
		}

		// instantiate new CameraModel with the camera reference and model details
		camera := CameraModel{
			camera:              eosCameraRef,
			szPortName:          C.GoString((*_Ctype_char)(&eosCameraDeviceInfo.szPortName[0])),
			szDeviceDescription: C.GoString((*_Ctype_char)(&eosCameraDeviceInfo.szDeviceDescription[0])),
			deviceSubType:       (uint32)(eosCameraDeviceInfo.deviceSubType),
			reserved:            (uint32)(eosCameraDeviceInfo.reserved),
		}
		cameras = append(cameras, camera)
	}

	return cameras, nil
}

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
