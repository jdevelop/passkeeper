package main

import (
	"fmt"
	"github.com/jdevelop/passkeeper/app/service/support"
	"github.com/jdevelop/passkeeper/device/rpi/display/oled"
	"log"
)

func main() {
	dev, err := oled.NewOLED(1, 60, 128, 64)
	if err != nil {
		log.Fatal(err)
	}

	db := support.MakeDashboard(dev, 5)
	for i := 0; i < 5; i++ {
		db.Log(fmt.Sprintf("%d", i))
	}
	fmt.Println("Done")
}
