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
	"unsafe"
)

type EOSClient struct{}

func NewEOSClient() *EOSClient {
	return &EOSClient{}
}

// Initialize a new EOSClient instance
func (e *EOSClient) Initialize() (err error) {
	eosError := C.EdsInitializeSDK()
	if eosError != C.EDS_ERR_OK {
		err = errors.New("Error when initializing Canon SDK")
	}
	return
}

// Release the EOSClient, must be called on termination
func (p *EOSClient) Release() (err error) {
	eosError := C.EdsTerminateSDK()
	if eosError != C.EDS_ERR_OK {
		err = errors.New("Error when terminating Canon SDK")
	}
	return
}

// Assign a callback function for when cameras are removed (disconnected).
func (e *EOSClient) SetCameraAddedHandler(h func()) {
	C.EdsSetCameraAddedHandler((*[0]byte)(unsafe.Pointer(&h)), nil)
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
		err := errors.New("Error when obtaining reference to list of cameras")
		return nil, err
	}
	defer C.EdsRelease((*C.struct___EdsObject)(unsafe.Pointer(&eosCameraList)))

	// get the number of cameras connected
	if eosError = C.EdsGetChildCount((*C.struct___EdsObject)(unsafe.Pointer(eosCameraList)), (*C.EdsUInt32)(unsafe.Pointer(&cameraCount))); eosError != C.EDS_ERR_OK {
		err := errors.New("Error when obtaining count of cameras")
		return nil, err
	}

	// get details for each camera detected
	cameras := make([]CameraModel, cameraCount)
	for i := 0; i < cameraCount; i++ {
		var eosCameraRef *C.EdsCameraRef
		if eosError = C.EdsGetChildAtIndex((*C.struct___EdsObject)(unsafe.Pointer(eosCameraList)), (C.EdsInt32)(i), (*C.EdsBaseRef)(unsafe.Pointer(&eosCameraRef))); eosError != C.EDS_ERR_OK {
			err := errors.New("Error when obtaining camera")
			return nil, err
		}

		var eosCameraDeviceInfo C.EdsDeviceInfo
		if eosError = C.EdsGetDeviceInfo((*C.struct___EdsObject)(unsafe.Pointer(eosCameraRef)), (*C.EdsDeviceInfo)(unsafe.Pointer(&eosCameraDeviceInfo))); eosError != C.EDS_ERR_OK {
			err := errors.New("Error when obtaining camera device info")
			return nil, err
		}

		// instantiate new CameraModel with the camera reference and model details
		cameras = append(cameras, CameraModel{
			camera:              eosCameraRef,
			szPortName:          C.GoString((*_Ctype_char)(&eosCameraDeviceInfo.szPortName[0])),
			szDeviceDescription: C.GoString((*_Ctype_char)(&eosCameraDeviceInfo.szDeviceDescription[0])),
			deviceSubType:       (uint32)(eosCameraDeviceInfo.deviceSubType),
			reserved:            (uint32)(eosCameraDeviceInfo.reserved),
		})
	}

	return cameras, nil
}
