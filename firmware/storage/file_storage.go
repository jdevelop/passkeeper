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

func readCredentials(s *PlainText) ([]firmware.Credentials, error) {
	s.file.Seek(0, 0)
	bytes, err := ioutil.ReadAll(s.file)
	if err != nil {
		return nil, err
	}

	if len(bytes) == 0 {
		return make([]firmware.Credentials, 0), nil
	}

	pt, err := decrypt([]byte(s.key), bytes)

	if err != nil {
		return nil, err
	}

	var seeds []firmware.Credentials

	err = json.Unmarshal(pt, &seeds)
	if err != nil {
		return nil, err
	}

	return seeds, nil
}

func (s *PlainText) ReadCredentials(id string) (*firmware.Credentials, error) {

	passwords, err := readCredentials(s)
	if err != nil {
		return nil, err
	}
	var password firmware.Credentials
	for _, s := range passwords {
		if s.Id == id {
			password = s
			break
		}
	}
	return &password, nil
}

func (s *PlainText) RemoveCredentials(id string) error {
	passwords, err := readCredentials(s)
	if err != nil {
		return err
	}
	var pos = -1
	for i, password := range passwords {
		if password.Id == id {
			pos = i
			break
		}
	}
	if pos > -1 {
		last := len(passwords) - 1
		passwords[pos] = passwords[last]
		return s.writeCredentialsToFile(passwords[:last])
	}
	return fmt.Errorf("cant find seed for %s", id)
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

func (s *PlainText) writeCredentialsToFile(seeds []firmware.Credentials) error {
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

func (s *PlainText) WriteCredentials(seed firmware.Credentials) error {

	seeds, err := readCredentials(s)

	if err != nil {
		return err
	}

	exI := -1

	for i, s := range seeds {
		if s.Id == seed.Id {
			exI = i
			break
		}
	}

	if exI == -1 {
		seeds = append(seeds, seed)
	} else {
		seeds[exI] = seed
	}

	return s.writeCredentialsToFile(seeds)

}

func (s *PlainText) ListCredentials() ([]firmware.Credentials, error) {
	_, err := s.file.Seek(0, 0)
	if err != nil {
		return nil, err
	}
	seeds, err := readCredentials(s)
	if err != nil {
		return nil, err
	}

	return seeds, nil
}

func (s *PlainText) Close() error {
	return s.file.Close()
}

var (
	p PlainText
	_ CredentialsStorageRead   = &p
	_ CredentialsStorageWrite  = &p
	_ CredentialsStorageList   = &p
	_ CredentialsStorageRemove = &p
)
