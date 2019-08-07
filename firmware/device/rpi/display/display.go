package display

import (
	"fmt"

	"github.com/jdevelop/passkeeper/firmware"
)

type HardwareDisplay interface {
	Cls()
	SetCursor(column, row uint8)
	Print(message string)
	Draw()
}

type TextBlockDisplay struct {
	realDisplay  HardwareDisplay
	content      *[]firmware.Credentials
	offset       int
	screenHeight int
}

func MakeTextDisplay(hw HardwareDisplay, content *[]firmware.Credentials, height int) *TextBlockDisplay {
	return &TextBlockDisplay{
		realDisplay:  hw,
		content:      content,
		offset:       0,
		screenHeight: height,
	}
}

func (d *TextBlockDisplay) Refresh() {
	cLen := len(*d.content)
	d.realDisplay.Cls()
	if d.content == nil && cLen == 0 {
		fmt.Println("Content empty!")
		d.realDisplay.SetCursor(0, 0)
		d.realDisplay.Print("No keys found")
		return
	}

	windowTop := d.offset
	windowBottom := d.offset + d.screenHeight - 1

	if windowBottom >= cLen {
		windowBottom = cLen - 1
		windowTop = cLen - d.screenHeight
	}

	if windowTop < 0 {
		windowTop = 0
	}

	for i := 0; windowTop+i <= windowBottom; i++ {
		d.realDisplay.SetCursor(uint8(i), 0)
		d.realDisplay.Print(" " + (*d.content)[windowTop+i].Service)
	}
	d.realDisplay.SetCursor(uint8(d.offset-windowTop), 0)
	d.realDisplay.Print(">")
	d.realDisplay.Draw()
}

func (d *TextBlockDisplay) ScrollUp(lines int) {
	d.offset = d.offset - lines
	if d.offset < 0 {
		d.offset = 0
	}
	d.Refresh()
}

func (d *TextBlockDisplay) ScrollDown(lines int) {
	d.offset = d.offset + lines
	if d.content != nil && d.offset >= len(*d.content) {
		d.offset = len(*d.content) - 1 // last item selected
	}
	if d.offset < 0 {
		d.offset = 0
	}
	d.Refresh()
}
