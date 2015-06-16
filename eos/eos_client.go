package eos

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework EDSDK
#define __MACOS__ 1
#include <EDSDK/EDSDK.h>
#include <stdlib.h>
*/
import "C"
import (
	"errors"
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
