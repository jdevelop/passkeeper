package storage_test

import (
	"fmt"
	"github.com/jdevelop/passkeeper"
	"github.com/jdevelop/passkeeper/storage"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
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
	assert.Equal(t, passkeeper.Seed{SeedId: "HELLO", SeedSecret: "WORLD"}, data)

	err = txt.SaveSeed(passkeeper.Seed{SeedId: "HELLO", SeedSecret: "pas"})
	assert.Nil(t, err)

	data, err = txt.LoadSeed("HELLO")
	fmt.Println(err)
	assert.Nil(t, err)
	assert.Equal(t, passkeeper.Seed{SeedId: "HELLO", SeedSecret: "pas"}, data)

	data, err = txt.LoadSeed("HELLOW")
	assert.Nil(t, err)
	assert.Equal(t, passkeeper.Seed{}, data)

	seeds, err := txt.ListSeeds()
	assert.Nil(t, err)
	assert.EqualValues(t, seeds, []string{"HELLO"})
	txt.Close()
}
