package main

import (
	"bufio"
	"fmt"
	"github.com/jdevelop/passkeeper/firmware"
	"github.com/jdevelop/passkeeper/firmware/app"
	"github.com/jdevelop/passkeeper/firmware/device/rpi"
	"github.com/jdevelop/passkeeper/firmware/pass"
	"github.com/jdevelop/passkeeper/firmware/storage"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"os"
	"periph.io/x/periph/host"
	"syscall"
)

var (
	cli = kingpin.New("util", "The utility application for seed management")

	mode = cli.Flag("mode", "Storage password source (card or manual entry)").
		Default("card").Enum("term", "card")

	config = cli.Flag("config", "Path to the config file").String()

	listCmd = cli.Command("list", "List available seeds")

	addCmd = cli.Command("add", "Add new seed")

	delCmd   = cli.Command("remove", "Remove seed by name")
	seedName = delCmd.Arg("id", "Seed id").Required().String()

	echoCmd    = cli.Command("echo", "Echo string to keyboard")
	echoString = echoCmd.Arg("text", "Text to echo").Required().String()

	resetCard     = cli.Command("reset_card", "Resets the card access")
	resetPassword = cli.Command("reset_password", "Resets the storage password")
	showPassword  = cli.Command("show_password", "Displays the current storage password")
)

func main() {

	cmdAlias, err := cli.Parse(os.Args[1:])

	if err != nil {
		log.Fatal(err)
	}

	var (
		cardPassword []byte
		cardAccess   [6]byte
	)

	c, err := app.LoadConfig(*config)
	if err != nil {
		log.Fatal(err)
	}

	readTerminalPwd := func() (pwd []byte, err error) {
		fmt.Print("Password >: ")
		pwd, err = terminal.ReadPassword(syscall.Stdin)
		return
	}

	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	openStorage := func() *storage.PlainText {

		switch *mode {
		case "term":
			pwdRaw, err := readTerminalPwd()
			if err != nil {
				log.Fatal(err)
			}
			cardPassword = storage.BuildKey(pwdRaw)
		case "card":
			rfid, err := pass.NewRFIDPass(c.Rfid.RfidAccessKey, c.Rfid.RfidAccessSector, c.Rfid.RfidAccessBlock)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Tap the card")
			cardPassword, err = rfid.GetCurrentPassword()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Password read successfully")
		}

		strg, err := storage.NewPlainText(c.Seeds.SeedFile, cardPassword)
		if err != nil {
			log.Fatal(err)
		}
		return strg
	}

	switch cmdAlias {
	case listCmd.FullCommand():
		strg := openStorage()
		seeds, err := strg.ListSeeds()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Seeds:")
		for _, seed := range seeds {
			fmt.Println("\t", seed)
		}
	case addCmd.FullCommand():
		strg := openStorage()
		fmt.Println("Seed name:")
		rdr := bufio.NewReader(os.Stdin)
		seedName, _, err := rdr.ReadLine()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Seed:")
		seed1, err := terminal.ReadPassword(syscall.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Confirm seed:")
		seed2, err := terminal.ReadPassword(syscall.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		if string(seed1) != string(seed2) {
			log.Fatalln("Seeds do not match, aborting")
		}
		err = strg.SaveSeed(passkeeper.Seed{
			SeedId:     string(seedName),
			SeedSecret: string(seed1),
		})
		if err != nil {
			log.Fatal(err)
		}
		err = strg.Close()
		if err != nil {
			log.Fatal(err)
		}

	case resetCard.FullCommand():
		pwdRaw, err := readTerminalPwd()
		if err != nil {
			log.Fatal(err)
		}
		copy(cardAccess[:], storage.BuildKey(pwdRaw)[:6])
		fmt.Println("Opening card")
		rfid, err := pass.NewRFIDPass(c.Rfid.RfidAccessKey, c.Rfid.RfidAccessSector, c.Rfid.RfidAccessBlock)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Tap the card to reset the card key")
		if err = rfid.ResetAccessKey(cardAccess, c.Rfid.RfidAccessSector); err != nil {
			log.Fatal(err)
		}
		copy(c.Rfid.RfidAccessKey[:], cardAccess[:])
		err = app.SaveConfig(*config, &c)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Access key updated successfully: %v\n", cardAccess)
		if err := rfid.Close(); err != nil {
			log.Fatal(err)
		}

	case resetPassword.FullCommand():
		pwdRaw, err := readTerminalPwd()
		if err != nil {
			log.Fatal(err)
		}

		cardPassword = storage.BuildKey(pwdRaw)

		fmt.Println("Initializing card")

		rfid, err := pass.NewRFIDPass(c.Rfid.RfidAccessKey, c.Rfid.RfidAccessSector, c.Rfid.RfidAccessBlock)
		if err != nil {
			log.Fatal(err)
		}
		err = rfid.ResetPassword(cardPassword, c.Rfid.RfidAccessSector, c.Rfid.RfidAccessBlock)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Password updated successfully")
		if err := rfid.Close(); err != nil {
			log.Fatal(err)
		}

	case showPassword.FullCommand():
		fmt.Printf("Initializing card: sector %d, block: %d\n", c.Rfid.RfidAccessSector, c.Rfid.RfidAccessBlock)
		rfid, err := pass.NewRFIDPass(c.Rfid.RfidAccessKey, c.Rfid.RfidAccessSector, c.Rfid.RfidAccessBlock)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Reading current password using sector %d, block %d, keys: %v\n", c.Rfid.RfidAccessSector, c.Rfid.RfidAccessBlock, c.Rfid.RfidAccessKey)
		pwdBytes, err := rfid.GetCurrentPassword()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Current password bytes are ", pwdBytes)
		if err := rfid.Close(); err != nil {
			log.Fatal(err)
		}

	case delCmd.FullCommand():
		strg := openStorage()
		if err := strg.RemoveSeed(*seedName); err != nil {
			log.Fatal(err)
		}

	case echoCmd.FullCommand():
		board, err := rpi.InitLinuxStack(rpi.StackParams{HasEthernet: true, HasSerial: true})
		if err != nil {
			log.Fatal(err)
		}
		if err := board.WriteString(*echoString); err != nil {
			log.Fatal(err)
		}
	default:
		cli.Usage(nil)
		return
	}

	fmt.Println("Done!")

}
