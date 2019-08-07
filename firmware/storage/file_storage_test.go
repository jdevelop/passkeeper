package storage

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/jdevelop/passkeeper/firmware"
	"github.com/stretchr/testify/assert"
)

func TestCredsGeneration(t *testing.T) {
	os.Remove("/tmp/passkeeperseed_plain_text.txt")
	txt, err := NewPlainText("/tmp/passkeeperseed_plain_text.txt", []byte("a very very very very secret key"))
	assert.Nil(t, err, "Error creating plaintext storage")
	data, err := txt.LoadCredentials("HELLO")
	assert.Nil(t, err)
	err = txt.SaveCredentials(firmware.Credentials{Service: "HELLO", Secret: "WORLD"})
	assert.Nil(t, err)
	data, err = txt.LoadCredentials("HELLO")
	assert.Nil(t, err)
	assert.Equal(t, firmware.Credentials{Service: "HELLO", Secret: "WORLD"}, *data)

	err = txt.SaveCredentials(firmware.Credentials{Service: "HELLO", Secret: "pas"})
	assert.Nil(t, err)

	data, err = txt.LoadCredentials("HELLO")
	fmt.Println(err)
	assert.Nil(t, err)
	assert.Equal(t, firmware.Credentials{Service: "HELLO", Secret: "pas"}, *data)

	data, err = txt.LoadCredentials("HELLOW")
	assert.Nil(t, err)
	assert.Equal(t, firmware.Credentials{}, *data)

	creds, err := txt.ListCredentialss()
	assert.Nil(t, err)
	assert.EqualValues(t, creds, []string{"HELLO"})
	txt.Close()
}

func TestBackupFile(t *testing.T) {
	const (
		filename = "/tmp/passkeeper.test"
		content  = "Oh my password"
	)
	testDataFile, err := os.Create(filename)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := testDataFile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	name, err := backupFile(testDataFile)
	if err != nil {
		t.Fatal(err)
	}
	data, err := ioutil.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != content {
		t.Fatalf("Expected '%s', got '%s' from %s", content, string(data), name)
	}
}
