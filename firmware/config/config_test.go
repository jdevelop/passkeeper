package config

import (
	"os"
	"reflect"
	"testing"
)

const confPath = "/tmp/goseed-test"

var newKey = [...]byte{0, 1, 0, 1, 0, 1}

func TestConfigLoad(t *testing.T) {
	f, err := os.Create(confPath)
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	if _, err := f.WriteString("{}"); err != nil {
		t.Error(err)
	}

	c, err := LoadConfig(confPath)
	if err != nil {
		t.Error(err)
	}

	c.Rfid.RfidAccessBlock = 100500
	c.Rfid.RfidAccessSector = 100501
	c.Rfid.RfidAccessKey = newKey
	if err = SaveConfig(confPath, &c); err != nil {
		t.Error(err)
	}

	c, err = LoadConfig(confPath)

	if c.Rfid.RfidAccessBlock != 100500 {
		t.Errorf("Expected access block 100500, was %d", c.Rfid.RfidAccessBlock)
	}

	if c.Rfid.RfidAccessSector != 100501 {
		t.Errorf("Expected access sector 100501, was %d", c.Rfid.RfidAccessSector)
	}

	if !reflect.DeepEqual(c.Rfid.RfidAccessKey, newKey) {
		t.Errorf("Keys are not equal: %v != %v", c.Rfid.RfidAccessKey, newKey)
	}
}
