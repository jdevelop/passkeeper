package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jdevelop/passkeeper/firmware"
)

type PlainText struct {
	filePath string
	key      []byte
}

func NewPlainText(filename string, key []byte) (*PlainText, error) {
	_, err := os.Stat(filename)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	keyData := make([]byte, len(key))
	copy(keyData, key)
	return &PlainText{
		filePath: filename,
		key:      keyData,
	}, nil
}

func readCredentials(s *PlainText) ([]firmware.Credentials, error) {
	bytes, err := ioutil.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return make([]firmware.Credentials, 0), nil
		}
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

func (s *PlainText) UpdateKey(newKey []byte) error {
	copy(s.key, newKey)
	return nil
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

func backupFile(src string) (string, error) {
	_, err := os.Stat(src)
	if err != nil && os.IsNotExist(err) {
		return "", nil
	}
	parentDir := filepath.Dir(src)
	currentFile := filepath.Base(src)
	newFile := fmt.Sprintf("%s.%d", filepath.Join(parentDir, currentFile), time.Now().Unix())
	dst, err := os.OpenFile(newFile, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return "", err
	}
	srcF, err := os.Open(src)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(dst, srcF)
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

	if name, err := backupFile(s.filePath); err != nil {
		return fmt.Errorf("Can't backup file '%s' : %v", name, err)
	}

	return ioutil.WriteFile(s.filePath, enc, 0600)
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
	seeds, err := readCredentials(s)
	if err != nil {
		return nil, err
	}

	return seeds, nil
}

func (s *PlainText) Close() error {
	return nil
}

func (s *PlainText) BackupStorage() (io.Reader, error) {
	data, err := ioutil.ReadFile(s.filePath)
	if err != nil {
		switch {
		case os.IsNotExist(err):
			return strings.NewReader("[]"), nil
		default:
			return nil, err
		}
	}
	pt, err := decrypt([]byte(s.key), data)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(pt), nil
}

func (s *PlainText) RestoreStorage(src io.Reader) error {
	data, err := ioutil.ReadAll(src)
	if err != nil {
		return err
	}
	var sample []firmware.Credentials
	if err := json.Unmarshal(data, &sample); err != nil {
		return err
	}
	if len(sample) == 0 {
		return fmt.Errorf("can't restore zero-length credentials list")
	}
	enc, err := encrypt(s.key, data)
	if err != nil {
		return err
	}
	_, err = backupFile(s.filePath)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(s.filePath, enc, 0600)
}

var (
	p PlainText
	_ CredentialsStorageRead   = &p
	_ CredentialsStorageWrite  = &p
	_ CredentialsStorageList   = &p
	_ CredentialsStorageRemove = &p
	_ CredentialsStorageBackup = &p
	_ CredentialsStorageBackup = &p
	_ MutableKey               = &p
)
