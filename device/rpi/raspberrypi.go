// +build linux

package rpi

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	gpio "github.com/jdevelop/gpio"
	rpio "github.com/jdevelop/gpio/rpi"
	"github.com/jdevelop/passkeeper/controls"
)

type led struct {
	standbyLedPin gpio.Pin
	errorLedPin   gpio.Pin
	busyLedPin    gpio.Pin
}

type RaspberryPi struct {
	Led          led
	InputControl controls.InputControl
}

type StackParams struct {
	HasSerial   bool
	HasEthernet bool
}

func (rpi *RaspberryPi) Close() {
	rpi.Led.busyLedPin.Close()
	rpi.Led.errorLedPin.Close()
	rpi.Led.standbyLedPin.Close()
}

func (rpi *RaspberryPi) Clear() {
	rpi.Led.busyLedPin.Clear()
	rpi.Led.errorLedPin.Clear()
	rpi.Led.standbyLedPin.Clear()
}

func (rpi *RaspberryPi) SelfCheckInprogress() (err error) {
	rpi.Clear()
	rpi.Led.busyLedPin.Set()
	return
}

func (rpi *RaspberryPi) SelfCheckComplete() (err error) {
	rpi.Clear()
	rpi.Led.standbyLedPin.Set()
	return
}

func (rpi *RaspberryPi) SelfCheckFailure(reason error) (err error) {
	rpi.Clear()
	rpi.Led.errorLedPin.Set()
	return
}

func (rpi *RaspberryPi) ReadyToTransmit() (err error) {
	rpi.Clear()
	rpi.Led.standbyLedPin.Set()
	return
}

func (rpi *RaspberryPi) TransmissionComplete() (err error) {
	rpi.Clear()
	rpi.Led.standbyLedPin.Set()
	return
}

func (rpi *RaspberryPi) TransmissionFailure(reason error) (err error) {
	fmt.Println("Transmission failed", reason)
	rpi.Clear()
	rpi.Led.errorLedPin.Set()
	return
}

type Led struct {
	standbyLedPin int
	errorLedPin   int
	busyLedPin    int
}

type Control struct {
	echoPin     int
	moveUpPin   int
	moveDownPin int
}

type DisplayControl struct {
}

type Board struct {
	Led     Led
	Control Control
}

func LedSettings(standbyLedPin, errorLedPin, busyPin int) Led {
	return Led{
		standbyLedPin: standbyLedPin,
		errorLedPin:   errorLedPin,
		busyLedPin:    busyPin,
	}
}

func ControlSettings(echoPin, upPin, downPin int) Control {
	return Control{
		echoPin:     echoPin,
		moveDownPin: downPin,
		moveUpPin:   upPin,
	}
}

func CreateBoard(settings Board) (rpi *RaspberryPi, err error) {
	busyLedPin, err := rpio.OpenPin(settings.Led.busyLedPin, gpio.ModeOutput)
	if err != nil {
		return
	}
	errorLedPin, err := rpio.OpenPin(settings.Led.errorLedPin, gpio.ModeOutput)
	if err != nil {
		return
	}
	standbyLedPin, err := rpio.OpenPin(settings.Led.standbyLedPin, gpio.ModeOutput)
	if err != nil {
		return
	}

	preparePin := func(pinNum int) (pin gpio.Pin, pinErr error) {
		fmt.Printf("Opening pin %d\n", pinNum)
		pin, pinErr = rpio.OpenPin(pinNum, gpio.ModeInput)
		if pinErr != nil {
			return
		}
		pin.PullUp()
		return
	}

	echoPin, err := preparePin(settings.Control.echoPin)
	if err != nil {
		return
	}
	upPin, err := preparePin(settings.Control.moveUpPin)
	if err != nil {
		return
	}
	downPin, err := preparePin(settings.Control.moveDownPin)
	if err != nil {
		return
	}

	rpi = &RaspberryPi{
		Led: led{
			standbyLedPin: standbyLedPin,
			errorLedPin:   errorLedPin,
			busyLedPin:    busyLedPin,
		},
	}

	err = echoPin.BeginWatch(gpio.EdgeFalling, func() {
		if rpi.InputControl != nil {
			rpi.InputControl.OnClickOk()
		}
	})

	if err != nil {
		return
	}

	err = upPin.BeginWatch(gpio.EdgeFalling, func() {
		if rpi.InputControl != nil {
			rpi.InputControl.OnClickUp()
		}
	})
	if err != nil {
		return
	}

	err = downPin.BeginWatch(gpio.EdgeFalling, func() {
		if rpi.InputControl != nil {
			rpi.InputControl.OnClickDown()
		}
	})
	if err != nil {
		return
	}

	return
}

func InitLinuxStack(params StackParams) (*VirtualKeyboard, error) {
	log("Init devices")
	_, err := os.Stat(deviceName)
	if err == nil || os.IsExist(err) {
		log("Keyboard exists, re-using descriptor")
		return &localKbd, nil
	}

	log("Create gadget %s", usbGadget)
	err = os.Mkdir(usbGadget, os.ModeDir)
	if err != nil {
		return nil, err
	}
	log("Set vendor data")
	if err = writeFiles(usbGadget, [][]string{
		{"0x1d6b", "idVendor"},
		{"0x0104", "idProduct"},
		{"0x0100", "bcdDevice"},
		{"0x0200", "bcdUSB"},
	}); err != nil {
		return nil, err
	}
	log("Create %s", usbStrings)
	if err = os.MkdirAll(usbStrings, os.ModeDir); err != nil {
		return nil, err
	}
	log("Set %s data", usbStrings)
	if err = writeFiles(usbStrings, [][]string{
		{"b65a0fe47231d98c", "serialnumber"},
		{"PASSKEEPER", "manufacturer"},
		{"PASSKEEPER HW", "product"},
	}); err != nil {
		return nil, err
	}
	log("Create config %s", usbConfig+"/strings/0x409")
	if err = os.MkdirAll(usbConfig+"/strings/0x409", os.ModeDir); err != nil {
		return nil, err
	}
	log("Create config %s", usbConfigStrings+"/configuration")
	if err = ioutil.WriteFile(usbConfigStrings+"/configuration", []byte("Config 1: ECM network"), filemode); err != nil {
		return nil, err
	}
	log("Setup power %s", usbConfig+"/MaxPower")
	if err = ioutil.WriteFile(usbConfig+"/MaxPower", []byte("200"), filemode); err != nil {
		return nil, err
	}

	if params.HasSerial {
		log("Setup serial %s", usbSerial)
		if err = os.MkdirAll(usbSerial, os.ModeDir); err != nil {
			return nil, err
		}
		log("Setup link %s -> %s", usbSerial, usbConfig+"/acm.usb0")
		if err = os.Symlink(usbSerial, usbConfig+"/acm.usb0"); err != nil {
			return nil, err
		}
	}

	if params.HasEthernet {
		log("Setting up ethernet")
		if err = ethernetUp(); err != nil {
			return nil, err
		}
	}

	log("Setting up usbHid %s", usbHid)
	if err = os.MkdirAll(usbHid, os.ModeDir); err != nil {
		return nil, err
	}

	if err = writeFiles(usbHid, [][]string{
		{"1", "protocol"},
		{"1", "subclass"},
		{"8", "report_length"},
	}); err != nil {
		return nil, err
	}

	log("Setting up usbHid %s", usbHid+"/report_desc")
	if err = ioutil.WriteFile(usbHid+"/report_desc", report[:], filemode); err != nil {
		return nil, err
	}

	log("Setting up link %s -> %s", usbHid, usbConfig+"/hid.usb0")
	if err = os.Symlink(usbHid, usbConfig+"/hid.usb0"); err != nil {
		return nil, err
	}

	files, err := ioutil.ReadDir("/sys/class/udc")
	if err != nil {
		return nil, err
	}

	var buffer string

	for _, file := range files {
		buffer = buffer + file.Name() + "\n"
	}

	log("Create UDC gadget %s", usbGadget+"/UDC")
	if err = ioutil.WriteFile(usbGadget+"/UDC", []byte(buffer), filemode); err != nil {
		return nil, err
	}

	time.Sleep(3 * time.Second)

	if params.HasEthernet {
		log("Setting up ethernet")
		if err = networkUp("usb0"); err != nil {
			return nil, err
		}

		log("Init DHCP")
		if _, err = dhcpUp("usb0"); err != nil {
			return nil, err
		}
	}

	return &localKbd, nil
}
