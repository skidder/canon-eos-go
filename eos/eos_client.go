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
