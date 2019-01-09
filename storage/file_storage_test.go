package storage

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/jdevelop/passkeeper"
	"github.com/stretchr/testify/assert"
)

func TestSeedGeneration(t *testing.T) {
	os.Remove("/tmp/passkeeperseed_plain_text.txt")
	txt, err := storage.NewPlainText("/tmp/passkeeperseed_plain_text.txt", []byte("a very very very very secret key"))
	assert.Nil(t, err, "Error creating plaintext storage")
	data, err := txt.LoadSeed("HELLO")
	assert.Nil(t, err)
	err = txt.SaveSeed(passkeeper.Seed{SeedId: "HELLO", SeedSecret: "WORLD"})
	assert.Nil(t, err)
	data, err = txt.LoadSeed("HELLO")
	assert.Nil(t, err)
	assert.Equal(t, passkeeper.Seed{SeedId: "HELLO", SeedSecret: "WORLD"}, *data)

	err = txt.SaveSeed(passkeeper.Seed{SeedId: "HELLO", SeedSecret: "pas"})
	assert.Nil(t, err)

	data, err = txt.LoadSeed("HELLO")
	fmt.Println(err)
	assert.Nil(t, err)
	assert.Equal(t, passkeeper.Seed{SeedId: "HELLO", SeedSecret: "pas"}, *data)

	data, err = txt.LoadSeed("HELLOW")
	assert.Nil(t, err)
	assert.Equal(t, passkeeper.Seed{}, *data)

	seeds, err := txt.ListSeeds()
	assert.Nil(t, err)
	assert.EqualValues(t, seeds, []string{"HELLO"})
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
