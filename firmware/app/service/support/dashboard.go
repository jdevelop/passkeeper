package support

import (
	"fmt"
	"github.com/jdevelop/passkeeper/firmware/device/rpi/display"
)

type dashboard struct {
	d      display.HardwareDisplay
	idx    int
	data   []string
	height int
}

func (d *dashboard) Log(msg string) {
	d.d.Cls()
	if d.idx == d.height {
		d.data = d.data[1:]
		d.data = append(d.data, msg)
	} else {
		d.data[d.idx] = msg
		d.idx = d.idx + 1
	}
	if d.d != nil {
		for i, m := range d.data {
			d.d.SetCursor(uint8(i), 0)
			d.d.Print(m)
		}
	}
	d.d.Draw()
	fmt.Println(msg)
}

func MakeDashboard(d display.HardwareDisplay, height int) (db *dashboard) {
	db = &dashboard{
		d:    d,
		data: make([]string, height),
		idx:  0,
	}
	return
}
