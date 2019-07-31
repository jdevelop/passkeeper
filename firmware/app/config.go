package app

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os/user"
	"time"
)

type leds struct {
	Error   int `json:"error"`
	Standby int `json:"standby"`
	Ready   int `json:"ready"`
}

type keys struct {
	Send int `json:"key_send"`
	Up   int `json:"key_up"`
	Down int `json:"key_down"`
}

type passwords struct {
	PasswordFile string `json:"file"`
}

type rfid struct {
	RfidAccessKey    [6]byte `json:"access_key,omitempty"`
	RfidAccessSector int     `json:"access_sector,omitempty"`
	RfidAccessBlock  int     `json:"access_block,omitempty"`
}

type lcd struct {
	DataPins []int `json:"data,omitempty"`
	RsPin    int   `json:"rs_pin,omitempty"`
	EPin     int   `json:"e_pin,omitempty"`
}

type oled struct {
	BusId  int `json:"bus_id,omitempty"`
	DevId  int `json:"dev_id,omitempty"`
	Width  int `json:"width,omitempty"`
	Height int `json:"height,omitempty"`
}

type Config struct {
	Leds      leds      `json:"leds"`
	Keys      keys      `json:"keys"`
	Passwords passwords `json:"passwords"`
	Rfid      rfid      `json:"rfid"`
	LCD       lcd       `json:"lcd"`
	OLED      oled      `json:"oled"`
}

var DefaultConfig = Config{
	Leds: leds{
		Error:   19,
		Standby: 20,
		Ready:   21,
	},
	Keys: keys{
		Send: 18,
		Down: 5,
		Up:   6,
	},
	Passwords: passwords{
		PasswordFile: "/root/passwordstorage.enc",
	},
	Rfid: rfid{
		RfidAccessBlock:  1,
		RfidAccessSector: 1,
		RfidAccessKey:    [...]byte{1, 2, 3, 4, 5, 6},
	},
	LCD: lcd{
		DataPins: []int{25, 24, 23, 22},
		EPin:     27,
		RsPin:    26,
	},
	OLED: oled{
		BusId:  1,
		DevId:  60,
		Width:  128,
		Height: 64,
	},
}

func resolveConfig(path string) (confPath string, err error) {
	if path == "" {
		u, errU := user.Current()
		if errU != nil {
			err = errU
			return
		}
		confPath = fmt.Sprintf("%s/.seedkeeper", u.HomeDir)
	} else {
		confPath = path
	}
	return
}

func LoadConfig(cfg string) (c Config, err error) {
	var jsonData []byte
	path, err := resolveConfig(cfg)
	if err != nil {
		return
	}
	jsonData, err = ioutil.ReadFile(path)
	if err != nil {
		return
	}

	c = DefaultConfig

	err = json.Unmarshal(jsonData, &c)

	return
}

func SaveConfig(cfg string, c *Config) (err error) {
	jsonData, err := json.Marshal(&c)
	if err != nil {
		return
	}
	path, err := resolveConfig(cfg)
	if err != nil {
		return
	}
	newName := fmt.Sprintf("%s.%d", path, time.Now().Unix())

	oldConf, err := ioutil.ReadFile(path)
	if err != nil {
		err = errors.Wrap(err, "Can't read the config file at "+path)
		return
	}
	err = ioutil.WriteFile(path, jsonData, 0600)
	if err != nil {
		err = errors.Wrap(err, "Can't write config file at "+path)
		return
	}

	err = ioutil.WriteFile(newName, oldConf, 0600)
	if err != nil {
		err = errors.Wrap(err, "Can't write backup config file at "+newName)
		return
	}
	return
}
