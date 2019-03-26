package main

import (
	"bytes"
	"image"
	"image/draw"
	"image/png"
	"log"

	"github.com/gobuffalo/packr"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/devices/ssd1306"
	"periph.io/x/periph/devices/ssd1306/image1bit"
	"periph.io/x/periph/host"
)

const (
	w = 128
	h = 64
)

var rect = image.Rect(0, 0, w, h)

func main() {
	// Load all the drivers:
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// Open a handle to the first available I²C bus:
	bus, err := i2creg.Open("")
	if err != nil {
		log.Fatal(err)
	}

	// Open a handle to a ssd1306 connected on the I²C bus:
	dev, err := ssd1306.NewI2C(bus, &ssd1306.Opts{
		W:       w,
		H:       h,
		Rotated: false,
	})
	if err != nil {
		log.Fatal(err)
	}

	box := packr.NewBox("./image")

	data, err := box.Find("splash.png")
	if err != nil {
		log.Fatal(err)
	}

	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}

	imgbw := image1bit.NewVerticalLSB(rect)
	draw.Draw(imgbw, rect, img, image.Point{}, draw.Src)

	err = dev.Draw(rect, imgbw, image.Point{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Done")

	bus.Close()

}
