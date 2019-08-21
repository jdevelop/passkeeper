package pass

import (
	"fmt"
	"time"

	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/conn/spi/spireg"
	"periph.io/x/periph/experimental/devices/mfrc522"
	"periph.io/x/periph/experimental/devices/mfrc522/commands"
)

type RFID struct {
	rfid                *mfrc522.Dev
	currentPass         string
	cardKey             [6]byte
	pwdSector, pwdBlock int
	currentPassOk       bool
}

var ZeroKey = [...]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

var defaultKey = [...]byte{0xca, 0xfe, 0xba, 0xbe, 0, 0}

const (
	//resetPinStr = "13"
	//irqPinStr   = "12"
	resetPinStr = "27"
	irqPinStr   = "17"
)

type config struct {
	resetPin, irqPin string
}

var defaultConfig = config{
	resetPin: "27",
	irqPin:   "17",
}

type PinConfF func(*config) *config

func WithResetPin(pin string) PinConfF {
	return func(c *config) *config {
		c.resetPin = pin
		return c
	}
}

func WithIRQPin(pin string) PinConfF {
	return func(c *config) *config {
		c.irqPin = pin
		return c
	}
}

func NewRFIDPass(cardKey [6]byte, pwdSector, pwdBlock int, confs ...PinConfF) (*RFID, error) {
	dev, err := spireg.Open("")
	if err != nil {
		return nil, err
	}

	c := defaultConfig
	{
		cc := &c
		for _, f := range confs {
			cc = f(cc)
		}
		c = *cc
	}

	resetPin := gpioreg.ByName(c.resetPin)
	if resetPin == nil {
		return nil, fmt.Errorf("can't open reset pin")
	}
	irqPin := gpioreg.ByName(c.irqPin)
	if irqPin == nil {
		return nil, fmt.Errorf("can't open irq pin")
	}

	nr, err := mfrc522.NewSPI(dev, resetPin, irqPin)
	if err != nil {
		return nil, err
	}

	return &RFID{
		rfid:      nr,
		cardKey:   cardKey,
		pwdBlock:  pwdBlock,
		pwdSector: pwdSector,
	}, nil
}

var (
	cardIRQTimeout = 30 * time.Second
)

func (r *RFID) GetCurrentPassword() ([]byte, error) {
	if err := r.rfid.LowLevel.WaitForEdge(cardIRQTimeout); err != nil {
		return nil, err
	}
	pwd, err := r.rfid.ReadCard(cardIRQTimeout, commands.PICC_AUTHENT1A, r.pwdSector, r.pwdBlock, r.cardKey)
	if err != nil {
		return nil, err
	}
	return pwd, nil
}

func (r *RFID) ResetAccessKey(newKeyArr [6]byte, sector int) error {
	if err := r.rfid.LowLevel.WaitForEdge(cardIRQTimeout); err != nil {
		return err
	}
	fmt.Printf("Card key %v => %v\n", r.cardKey, newKeyArr)
	return r.rfid.WriteSectorTrail(cardIRQTimeout, commands.PICC_AUTHENT1A, sector, newKeyArr, newKeyArr,
		&mfrc522.BlocksAccess{
			B0: mfrc522.AnyKeyRWID,
			B1: mfrc522.AnyKeyRWID,
			B2: mfrc522.AnyKeyRWID,
			B3: mfrc522.KeyA_RN_WA_BITS_RA_WA_KeyB_RA_WA,
		}, r.cardKey)
}

func (r *RFID) ResetPassword(newPassword []byte, sector, block int) error {
	if len(newPassword) != 16 {
		return fmt.Errorf("Password length must be of size 16 - found %d", len(newPassword))
	}

	var pwdArr [16]byte

	for i, v := range newPassword {
		pwdArr[i] = v
	}

	if err := r.rfid.LowLevel.WaitForEdge(cardIRQTimeout); err != nil {
		return err
	}
	return r.rfid.WriteCard(cardIRQTimeout, commands.PICC_AUTHENT1A, sector, block, pwdArr, r.cardKey)
}

func (r *RFID) GetCardKey() [6]byte {
	return r.cardKey
}

func (r *RFID) Close() error {
	return r.rfid.LowLevel.StopCrypto()
}
