package lcd

import (
	"fmt"
	lcd "github.com/jdevelop/golang-rpi-extras/lcd_hd44780"
	"github.com/jdevelop/passkeeper/device/rpi/display"
)

type Display struct {
	*lcd.PiLCD4
}

func MakeLCDDisplay(dataPins []int, ePin, rsPin int) (d display.HardwareDisplay, err error) {
	disp, err := lcd.NewLCD4(dataPins, rsPin, ePin)
	fmt.Printf("Init display: data %v, e: %v, rs: %v\n", dataPins, ePin, rsPin)
	disp.Init()
	d = &Display{
		PiLCD4: &disp,
	}
	return
}

func (d *Display) Draw() {
	// do nothing
}
