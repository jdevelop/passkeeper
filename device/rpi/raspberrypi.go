// +build linux

package rpi

import (
	"io/ioutil"
	"os"
	"time"

	gpio "github.com/jdevelop/gpio"
	rpio "github.com/jdevelop/gpio/rpi"
	"github.com/jdevelop/passkeeper/controls"
)

type (
	led struct {
		standbyLedPin gpio.Pin
		errorLedPin   gpio.Pin
		busyLedPin    gpio.Pin
	}

	RaspberryPi struct {
		Led          led
		InputControl controls.InputControl
	}

	StackParams struct {
		HasSerial   bool
		HasEthernet bool
	}

	Led struct {
		standbyLedPin int
		errorLedPin   int
		busyLedPin    int
	}

	Control struct {
		echoPin     int
		moveUpPin   int
		moveDownPin int
	}

	DisplayControl struct {
	}

	Board struct {
		Led     Led
		Control Control
	}
)

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

func (rpi *RaspberryPi) SelfCheckInprogress() error {
	rpi.Clear()
	rpi.Led.busyLedPin.Set()
	return nil
}

func (rpi *RaspberryPi) SelfCheckComplete() error {
	rpi.Clear()
	rpi.Led.standbyLedPin.Set()
	return nil
}

func (rpi *RaspberryPi) SelfCheckFailure(reason error) error {
	rpi.Clear()
	rpi.Led.errorLedPin.Set()
	return nil
}

func (rpi *RaspberryPi) ReadyToTransmit() error {
	rpi.Clear()
	rpi.Led.standbyLedPin.Set()
	return nil
}

func (rpi *RaspberryPi) TransmissionComplete() error {
	rpi.Clear()
	rpi.Led.standbyLedPin.Set()
	return nil
}

func (rpi *RaspberryPi) TransmissionFailure(reason error) error {
	log("Transmission failed: %v", reason)
	rpi.Clear()
	rpi.Led.errorLedPin.Set()
	return nil
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

func CreateBoard(settings Board) (*RaspberryPi, error) {
	busyLedPin, err := rpio.OpenPin(settings.Led.busyLedPin, gpio.ModeOutput)
	if err != nil {
		return nil, err
	}
	errorLedPin, err := rpio.OpenPin(settings.Led.errorLedPin, gpio.ModeOutput)
	if err != nil {
		return nil, err
	}
	standbyLedPin, err := rpio.OpenPin(settings.Led.standbyLedPin, gpio.ModeOutput)
	if err != nil {
		return nil, err
	}

	preparePin := func(pinNum int) (pin gpio.Pin, pinErr error) {
		log("Opening pin %d\n", pinNum)
		pin, pinErr = rpio.OpenPin(pinNum, gpio.ModeInput)
		if pinErr != nil {
			return nil, err
		}
		pin.PullUp()
		return nil, err
	}

	echoPin, err := preparePin(settings.Control.echoPin)
	if err != nil {
		return nil, err
	}
	upPin, err := preparePin(settings.Control.moveUpPin)
	if err != nil {
		return nil, err
	}
	downPin, err := preparePin(settings.Control.moveDownPin)
	if err != nil {
		return nil, err
	}

	rpi := &RaspberryPi{
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
		return nil, err
	}

	err = upPin.BeginWatch(gpio.EdgeFalling, func() {
		if rpi.InputControl != nil {
			rpi.InputControl.OnClickUp()
		}
	})
	if err != nil {
		return nil, err
	}

	err = downPin.BeginWatch(gpio.EdgeFalling, func() {
		if rpi.InputControl != nil {
			rpi.InputControl.OnClickDown()
		}
	})
	if err != nil {
		return nil, err
	}

	return rpi, nil
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
