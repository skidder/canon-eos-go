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

func (e *EOSClient) Initialize() (err error) {
	eosError := C.EdsInitializeSDK()
	if eosError != C.EDS_ERR_OK {
		err = errors.New("Error when initializing Canon SDK")
	}
	return
}

func (p *EOSClient) Close() (err error) {
	eosError := C.EdsTerminateSDK()
	if eosError != C.EDS_ERR_OK {
		err = errors.New("Error when terminating Canon SDK")
	}
	return
}

func (e *EOSClient) GetCameraModels() (cameras []CameraModel, err error) {
	var eosCameraList *C.EdsCameraListRef
	var eosError C.EdsError
	var cameraCount int

	// get a reference to the cameras list record
	if eosError = C.EdsGetCameraList((*C.EdsCameraListRef)(unsafe.Pointer(&eosCameraList))); eosError != C.EDS_ERR_OK {
		err = errors.New("Error when obtaining reference to list of cameras")
		return
	}
	defer C.EdsRelease((*C.struct___EdsObject)(unsafe.Pointer(&eosCameraList)))

	// get the number of cameras connected
	if eosError = C.EdsGetChildCount((*C.struct___EdsObject)(unsafe.Pointer(eosCameraList)), (*C.EdsUInt32)(unsafe.Pointer(&cameraCount))); eosError != C.EDS_ERR_OK {
		err = errors.New("Error when obtaining count of cameras")
		return
	}

	// get details for each camera detected
	for i := 0; i < cameraCount; i++ {
		var eosCameraRef *C.EdsCameraRef
		if eosError = C.EdsGetChildAtIndex((*C.struct___EdsObject)(unsafe.Pointer(eosCameraList)), (C.EdsInt32)(i), (*C.EdsBaseRef)(unsafe.Pointer(&eosCameraRef))); eosError != C.EDS_ERR_OK {
			err = errors.New("Error when obtaining camera")
			return
		}

		var eosCameraDeviceInfo C.EdsDeviceInfo
		if eosError = C.EdsGetDeviceInfo((*C.struct___EdsObject)(unsafe.Pointer(eosCameraRef)), (*C.EdsDeviceInfo)(unsafe.Pointer(&eosCameraDeviceInfo))); eosError != C.EDS_ERR_OK {
			err = errors.New("Error when obtaining camera device info")
			return
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

	return
}
