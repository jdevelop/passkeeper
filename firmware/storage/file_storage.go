package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/jdevelop/passkeeper/firmware"
)

type PlainText struct {
	file *os.File
	key  []byte
}

func NewPlainText(filename string, key []byte) (*PlainText, error) {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	return &PlainText{
		file: f,
		key:  key,
	}, nil
}

func readSeeds(s *PlainText) ([]passkeeper.Seed, error) {
	s.file.Seek(0, 0)
	bytes, err := ioutil.ReadAll(s.file)
	if err != nil {
		return nil, err
	}

	if len(bytes) == 0 {
		return make([]passkeeper.Seed, 0), nil
	}

	pt, err := decrypt([]byte(s.key), bytes)

	if err != nil {
		return nil, err
	}

	var seeds []passkeeper.Seed

	err = json.Unmarshal(pt, &seeds)
	if err != nil {
		return nil, err
	}

	return seeds, nil
}

func (s *PlainText) LoadSeed(id string) (*passkeeper.Seed, error) {

	seeds, err := readSeeds(s)
	if err != nil {
		return nil, err
	}
	var seed passkeeper.Seed
	for _, s := range seeds {
		if s.SeedId == id {
			seed = s
			break
		}
	}
	return &seed, nil
}

func (s *PlainText) RemoveSeed(key string) error {
	seeds, err := readSeeds(s)
	if err != nil {
		return err
	}
	var pos = -1
	for i, seed := range seeds {
		if seed.SeedId == key {
			pos = i
			break
		}
	}
	if pos > -1 {
		last := len(seeds) - 1
		seeds[pos] = seeds[len(seeds)-1]
		return s.writeSeedsToFile(seeds[:last])
	}
	return fmt.Errorf("cant find seed for %s", key)
}

func backupFile(src *os.File) (_ string, err error) {
	if _, err := src.Seek(0, 0); err != nil {
		return "", err
	}
	parentDir := filepath.Dir(src.Name())
	currentFile := filepath.Base(src.Name())
	newFile := fmt.Sprintf("%s.%d", filepath.Join(parentDir, currentFile), time.Now().Unix())
	dst, err := os.OpenFile(newFile, os.O_CREATE|os.O_RDWR, 0600)
	defer func() {
		if err == nil {
			err = dst.Close()
		} else {
			dst.Close()
		}
	}()
	if err != nil {
		return "", err
	}
	_, err = io.Copy(dst, src)
	return newFile, err
}

func (s *PlainText) writeSeedsToFile(seeds []passkeeper.Seed) error {
	bytes, err := json.Marshal(seeds)
	if err != nil {
		return err
	}

	enc, err := encrypt([]byte(s.key), bytes)
	if err != nil {
		return err
	}

	if name, err := backupFile(s.file); err != nil {
		return fmt.Errorf("Can't backup file '%s' : %v", name, err)
	}

	if err := s.file.Truncate(0); err != nil {
		return err
	}
	_, err = s.file.WriteAt(enc, 0)
	if err != nil {
		return err
	}

	return s.file.Sync()
}

func (s *PlainText) SaveSeed(seed passkeeper.Seed) error {

	seeds, err := readSeeds(s)

	if err != nil {
		return err
	}

	exI := -1

	for i, s := range seeds {
		if s.SeedId == seed.SeedId {
			exI = i
			break
		}
	}

	if exI == -1 {
		seeds = append(seeds, seed)
	} else {
		seeds[exI] = seed
	}

	return s.writeSeedsToFile(seeds)

}

func (s *PlainText) ListSeeds() ([]string, error) {
	_, err := s.file.Seek(0, 0)
	if err != nil {
		return nil, err
	}
	seeds, err := readSeeds(s)
	if err != nil {
		return nil, err
	}

	res := make([]string, len(seeds))

	for i, s := range seeds {
		res[i] = s.SeedId
	}

	return res, nil
}

func (s *PlainText) Close() error {
	return s.file.Close()
}
