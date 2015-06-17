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
	"unsafe"
)

// Represents a connected camera model.  Must be released by invoking 'Release' function.
type CameraModel struct {
	camera *C.EdsCameraRef

	szPortName          string
	szDeviceDescription string
	deviceSubType       uint32
	reserved            uint32
}

// Releases reference to the camera
func (c *CameraModel) Release() {
	C.EdsRelease((*C.struct___EdsObject)(unsafe.Pointer(&c.camera)))
}
