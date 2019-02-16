package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jdevelop/passkeeper/app"
	"github.com/jdevelop/passkeeper/app/service/support"
	"github.com/jdevelop/passkeeper/controls/status"
	"github.com/jdevelop/passkeeper/device/rpi"
	"github.com/jdevelop/passkeeper/device/rpi/display"
	"github.com/jdevelop/passkeeper/device/rpi/display/lcd"
	"github.com/jdevelop/passkeeper/device/rpi/display/oled"
	"github.com/jdevelop/passkeeper/pass"
	"github.com/jdevelop/passkeeper/rest"
	"github.com/jdevelop/passkeeper/storage"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	"periph.io/x/periph/host"
)

var (
	cli         = kingpin.New("service", "Passkeeper")
	config      = cli.Flag("config", "Config path").Short('c').String()
	displayType = cli.Flag("display", "Display").Short('d').Default("lcd").Enum("lcd", "oled")
)

var currentSeedID = 0

var currentSeeds []string

func main() {

	_, err := cli.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	c, err := app.LoadConfig(*config)

	if err != nil {
		log.Fatal(err)
	}

	if _, err = host.Init(); err != nil {
		log.Fatal(err)
	}

	var hwDisplay display.HardwareDisplay
	var lines int

	switch *displayType {
	case "lcd":
		hwDisplay, err = lcd.MakeLCDDisplay(c.LCD.DataPins, c.LCD.EPin, c.LCD.RsPin)
		if err != nil {
			log.Fatal(err)
		}
		lines = 2
	case "oled":
		dsp, err := oled.NewOLED(c.OLED.BusId, c.OLED.DevId, c.OLED.Width, c.OLED.Height)
		hwDisplay = dsp
		if err != nil {
			log.Fatal(err)
		}
		lines = 5
	}

	db := support.MakeDashboard(hwDisplay, lines)

	db.Log("Init...")

	raspberry, err := rpi.CreateBoard(rpi.Board{
		Led:     rpi.LedSettings(c.Leds.Standby, c.Leds.Error, c.Leds.Ready),
		Control: rpi.ControlSettings(c.Keys.Send, c.Keys.Up, c.Keys.Down),
	})

	db.Log("Selfcheck")
	raspberry.SelfCheckInprogress()

	if err != nil {
		db.Log("Hardware failure")
		log.Fatal(err)
	}

	drv, err := rpi.InitLinuxStack(rpi.StackParams{HasSerial: true, HasEthernet: true})
	if err != nil {
		log.Fatal(err)
		db.Log("Kbd/net failure")
		raspberry.SelfCheckFailure(err)
	}

	rfid, err := pass.NewRFIDPass(c.Rfid.RfidAccessKey, c.Rfid.RfidAccessSector, c.Rfid.RfidAccessBlock)
	if err != nil {
		db.Log("RFID failure")
		raspberry.SelfCheckFailure(err)
		log.Fatal(err)
	}

	var pass []byte

	for {
		db.Log("Tap the RFID")
		pass, err = GetCurrentPassword(rfid, raspberry)

		if err != nil {
			db.Log("Key failed, retry")
			raspberry.SelfCheckFailure(err)
			fmt.Println("Error reading card, retrying", err)
			time.Sleep(1 * time.Second)
			raspberry.SelfCheckInprogress()
		} else {
			break
		}
	}

	db.Log("Storage open")

	raspberry.ReadyToTransmit()

	seedStorage, err := getStorage(c.Seeds.SeedFile, pass)

	if err != nil {
		db.Log("Storage failed")
		log.Fatal(err)
	}

	currentSeeds, err = seedStorage.ListSeeds()
	if err != nil {
		db.Log("Seed failed")
		log.Fatal(err)
	}

	lastCall := time.Now()
	lastUp := time.Now()
	lastDown := time.Now()

	checkTime := func(timeVar *time.Time) bool {
		if time.Since(*timeVar) < 500*time.Millisecond {
			return false
		}
		*timeVar = time.Now()
		return true
	}

	textDisplay := display.MakeTextDisplay(hwDisplay, &currentSeeds, lines)
	if err != nil {
		log.Fatal(err)
	}

	cf := support.ControlsStub{
		SendFunc: func() {
			if !checkTime(&lastCall) {
				return
			}
			raspberry.ReadyToTransmit()
			seeds, err := seedStorage.ListSeeds()
			if err != nil {
				return
			}
			seedSize := len(seeds)
			if seedSize == 0 {
				err = errors.New("No seed")
				db.Log("No seed")
				raspberry.TransmissionFailure(err)
				return
			}
			if currentSeedID >= seedSize {
				currentSeedID = seedSize - 1
			}
			if currentSeedID < 0 {
				currentSeedID = 0
			}
			seed, err := seedStorage.LoadSeed(seeds[currentSeedID])
			if err != nil {
				db.Log("Storage failed")
				raspberry.TransmissionFailure(err)
				return
			}
			err = drv.WriteString(seed.SeedSecret)
			if err != nil {
				db.Log("Keyboard failed")
				raspberry.TransmissionFailure(err)
				return
			}
			err = raspberry.TransmissionComplete()
			return
		},
		UpFunc: func() {
			if !checkTime(&lastUp) {
				return
			}
			fmt.Println("Up")
			currentSeedID = currentSeedID - 1
			textDisplay.ScrollUp(1)
		},
		DownFunc: func() {
			if !checkTime(&lastDown) {
				return
			}
			fmt.Println("Down")
			currentSeedID = currentSeedID + 1
			textDisplay.ScrollDown(1)
		},
	}

	raspberry.InputControl = &cf

	db.Log("Web started")

	time.AfterFunc(2*time.Second, func() {
		textDisplay.Refresh()
	})

	rest.Start("0.0.0.0", 80, seedStorage, func() {
		currentSeedID = 0
		currentSeeds, err = seedStorage.ListSeeds()
		textDisplay.Refresh()
	})
}

func GetCurrentPassword(provider pass.PasswordProvider, board status.StatusControl) ([]byte, error) {
	pwd, err := provider.GetCurrentPassword()
	if err != nil {
		board.SelfCheckFailure(err)
		return nil, err
	}
	return pwd, nil
}

func getStorage(filename string, pwd []byte) (*storage.PlainText, error) {
	return storage.NewPlainText(filename, pwd)
}
