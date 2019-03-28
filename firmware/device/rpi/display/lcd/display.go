package lcd

import (
	"fmt"

	"github.com/jdevelop/passkeeper/firmware/device/rpi/display"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/experimental/devices/hd44780"
)

type Display struct {
	dev *hd44780.Dev
}

func (d *Display) Cls() {
	d.dev.Reset()
}

func (d *Display) SetCursor(column uint8, row uint8) {
	d.dev.SetCursor(column, row)
}

func (d *Display) Print(message string) {
	d.dev.Print(message)
}

func (d *Display) Draw() {
}

func outPin(pin int) gpio.PinOut {
	return gpioreg.ByName(fmt.Sprintf("%d", pin))
}

func MakeLCDDisplay(dataPins []int, ePin, rsPin int) (*Display, error) {
	data := make([]gpio.PinOut, len(dataPins))
	for i := range dataPins {
		out := outPin(dataPins[i])
		if out == nil {
			return nil, wrapf("can't create out pin %d")
		}
	}
	ePinOut := outPin(ePin)
	if ePinOut == nil {
		return nil, wrapf("can't create e-pin %d", ePin)
	}
	rsPinOut := outPin(rsPin)
	if rsPinOut == nil {
		return nil, wrapf("can't create rs-pin %d", ePin)
	}
	disp, err := hd44780.New(data, rsPinOut, ePinOut)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Init display: data %v, e: %v, rs: %v\n", dataPins, ePin, rsPin)
	if err := disp.Reset(); err != nil {
		return nil, wrapf("can't init LCD display %v", err)
	}
	return &Display{
		disp,
	}, nil
}

var _ display.HardwareDisplay = &Display{}

func wrapf(msg string, args ...interface{}) error {
	return fmt.Errorf(msg, args...)
}
