package oled

import (
	"image"
	"image/color"

	"github.com/fogleman/gg"
	"golang.org/x/image/font/basicfont"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/devices/ssd1306"
)

var defaultFace = basicfont.Face7x13

type OLED struct {
	*ssd1306.Dev
	textCol, textRow int
	ctx              *gg.Context
}

var (
	defaultRect = image.Rect(0, 0, 128, 64)
	upperLeft   = image.Point{}
)

func NewOLED(bus, device, width, height int) (o *OLED, err error) {

	dev, err := i2creg.Open("")
	if err != nil {
		return
	}

	oled, err := ssd1306.NewI2C(dev, width, height, false)
	if err != nil {
		return
	}

	ctx := gg.NewContext(width, height)
	ctx.SetColor(color.Black)
	ctx.Clear()
	ctx.SetColor(color.White)

	oled.Draw(defaultRect, ctx.Image(), upperLeft)

	if err != nil {
		return
	}
	o = &OLED{
		Dev:     oled,
		textCol: 0,
		textRow: 0,
		ctx:     ctx,
	}
	return
}

func (d *OLED) Cls() {
	d.ctx.Push()
	d.ctx.SetColor(color.Black)
	d.ctx.Clear()
	d.ctx.Pop()
	d.Dev.Draw(defaultRect, d.ctx.Image(), upperLeft)
}

func (d *OLED) SetCursor(row, column uint8) {
	d.textCol = int(column)
	d.textRow = int(row)
}

func (d *OLED) Print(message string) {
	newX := float64(defaultFace.Advance * d.textCol)
	newY := float64((defaultFace.Height - 1) * (d.textRow + 1))
	d.ctx.DrawString(message, newX, newY)
}

func (d *OLED) Draw() {
	d.ctx.Stroke()
	d.Dev.Draw(defaultRect, d.ctx.Image(), upperLeft)
}

func (d *OLED) Close() {
	// Nothing to do here.
}
