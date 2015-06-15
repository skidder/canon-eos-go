package eos

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework EDSDK
#define __MACOS__ 1
#include <EDSDK/EDSDK.h>
#include <stdlib.h>
*/
import "C"
import ()

type EOSClient struct{}

func NewEOSClient() *EOSClient {
	return &EOSClient{}
}

func (e *EOSClient) Initialize() {
	C.EdsInitializeSDK()
}

func (p *EOSClient) Close() {
	defer C.EdsTerminateSDK()
}
