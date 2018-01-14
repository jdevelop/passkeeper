// +build linux

package rpi

import (
	"fmt"
	"io/ioutil"
	"os"
)

const (
	usbGadget = "/sys/kernel/config/usb_gadget/passkeeper"

	usbStrings       = usbGadget + "/strings/0x409"
	usbConfig        = usbGadget + "/configs/c.1"
	usbConfigStrings = usbConfig + "/strings/0x409"
	usbSerial        = usbGadget + "/functions/acm.usb0"
	usbHid           = usbGadget + "/functions/hid.usb0"
	name             = "passkeeper"
	deviceName       = "/dev/hidg0"

	filemode = 0600

	charOffset  = 0
	digitOffset = 26

	shiftBitsOffset = 0
	symOffset       = 2
)

var report = [...]byte{0x05, 0x01, 0x09, 0x06, 0xa1, 0x01, 0x05, 0x07, 0x19, 0xe0, 0x29, 0xe7, 0x15, 0x00, 0x25, 0x01, 0x75, 0x01, 0x95, 0x08, 0x81, 0x02, 0x95, 0x01, 0x75, 0x08, 0x81, 0x03, 0x95, 0x05, 0x75, 0x01, 0x05, 0x08, 0x19, 0x01, 0x29, 0x05, 0x91, 0x02, 0x95, 0x01, 0x75, 0x03, 0x91, 0x03, 0x95, 0x06, 0x75, 0x08, 0x15, 0x00, 0x25, 0x65, 0x05, 0x07, 0x19, 0x00, 0x29, 0x65, 0x81, 0x00, 0xc0}

func writeFiles(base string, data [][]string) error {
	createAndWrite := func(path string, data string) error {
		return ioutil.WriteFile(base+"/"+path, []byte(data), filemode)
	}
	for _, p := range data {
		if err := createAndWrite(p[1], p[0]); err != nil {
			return err
		}
	}
	return nil
}

type VirtualKeyboard struct{}

var localKbd = VirtualKeyboard{}

var debug = true

func log(msg string, params ...interface{}) {
	if debug {
		fmt.Printf(msg+"\n", params)
	}
}

func (ci *VirtualKeyboard) WriteString(content string) error {
	f, err := os.OpenFile(deviceName, os.O_RDWR, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

	stroke := make([]byte, 8)
	for _, ch := range content {
		key, err := ResolveScanKey(ch)
		if err != nil {
			continue
		}
		stroke[shiftBitsOffset] = key[1]
		stroke[symOffset] = key[0]
		if _, err := f.Write(stroke); err != nil {
			return err
		}
		stroke[shiftBitsOffset] = 0
		stroke[symOffset] = 0
		if _, err = f.Write(stroke); err != nil {
			return err
		}
	}
	return nil
}
