// +build linux

package rpi

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/jdevelop/passkeeper/firmware/controls/input"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
)

type (
	led struct {
		standbyLedPin gpio.PinOut
		errorLedPin   gpio.PinOut
		busyLedPin    gpio.PinOut
	}

	RaspberryPi struct {
		Led          led
		InputControl input.InputControl
		mainloop     chan struct{}
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

	Board struct {
		Led     Led
		Control Control
	}
)

func (rpi *RaspberryPi) String() string {
	return "BOARD"
}

func (rpi *RaspberryPi) Close() {
	rpi.Led.busyLedPin.Halt()
	rpi.Led.errorLedPin.Halt()
	rpi.Led.standbyLedPin.Halt()
}

func (rpi *RaspberryPi) Clear() {
	rpi.Led.busyLedPin.Out(gpio.Low)
	rpi.Led.errorLedPin.Out(gpio.Low)
	rpi.Led.standbyLedPin.Out(gpio.Low)
}

func (rpi *RaspberryPi) SelfCheckInprogress() error {
	rpi.Clear()
	rpi.Led.busyLedPin.Out(gpio.High)
	return nil
}

func (rpi *RaspberryPi) SelfCheckComplete() error {
	rpi.Clear()
	rpi.Led.standbyLedPin.Out(gpio.High)
	return nil
}

func (rpi *RaspberryPi) SelfCheckFailure(reason error) error {
	rpi.Clear()
	rpi.Led.errorLedPin.Out(gpio.High)
	return nil
}

func (rpi *RaspberryPi) ReadyToTransmit() error {
	rpi.Clear()
	rpi.Led.standbyLedPin.Out(gpio.High)
	return nil
}

func (rpi *RaspberryPi) TransmissionComplete() error {
	rpi.Clear()
	rpi.Led.standbyLedPin.Out(gpio.High)
	return nil
}

func (rpi *RaspberryPi) TransmissionFailure(reason error) error {
	log("Transmission failed: %v", reason)
	rpi.Clear()
	rpi.Led.errorLedPin.Out(gpio.High)
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

func bcmpin(pin int) string {
	return fmt.Sprintf("%d", pin)
}

func CreateBoard(settings Board) (*RaspberryPi, error) {
	busyLedPin := gpioreg.ByName(bcmpin(settings.Led.busyLedPin))
	if busyLedPin == nil {
		return nil, wrapf("can't find busy pin led")
	}
	errorLedPin := gpioreg.ByName(bcmpin(settings.Led.errorLedPin))
	if errorLedPin == nil {
		return nil, wrapf("can't find error pin led")
	}
	standbyLedPin := gpioreg.ByName(bcmpin(settings.Led.standbyLedPin))
	if standbyLedPin == nil {
		return nil, wrapf("can't find standby pin led")
	}

	preparePin := func(pinNum int) (gpio.PinIn, error) {
		log("Opening pin %d\n", pinNum)
		pin := gpioreg.ByName(bcmpin(pinNum))
		if pin == nil {
			return nil, wrapf("can't find input pin %d", pinNum)
		}
		if err := pin.In(gpio.PullUp, gpio.FallingEdge); err != nil {
			return nil, err
		}
		return pin, nil
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

	log("Board settings: %s", rpi.String())

	mainloop := make(chan struct{})

	watcher := func(p gpio.PinIn, f func()) {
		go func() {
			for {
				select {
				case _, open := <-mainloop:
					if !open {
						return
					}
				default:
					if p.WaitForEdge(time.Millisecond * 20) {
						if rpi.InputControl != nil {
							f()
						}
					}

				}
			}
		}()
	}

	watcher(echoPin, func() { rpi.InputControl.OnClickOk() })
	watcher(upPin, func() { rpi.InputControl.OnClickUp() })
	watcher(downPin, func() { rpi.InputControl.OnClickDown() })

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

	//time.Sleep(1 * time.Second)

	if params.HasEthernet {
		log("Setting up ethernet")
		if err = networkUp("usb0"); err != nil {
			return nil, err
		}

		log("Init DHCP")
		go func() {
			err := dhcpUp("usb0", localIPStr, leaseStartStr)
			if err != nil {
				log("Can't start DHCP server %v", err)
			}
		}()
	}

	return &localKbd, nil
}

func wrapf(err string, args ...interface{}) error {
	return fmt.Errorf(err, args...)
}
