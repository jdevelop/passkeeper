package storage

import (
	"errors"
	"io"

	"github.com/jdevelop/passkeeper/firmware"
)

var ZeroLengthPasswordList = errors.New("can't restore zero-length passwords")

type CredentialsStorageRead interface {
	ReadCredentials(string) (*firmware.Credentials, error)
}

type CredentialsStorageWrite interface {
	WriteCredentials(firmware.Credentials) error
}

type CredentialsStorageList interface {
	ListCredentials() ([]firmware.Credentials, error)
}

type CredentialsStorageRemove interface {
	RemoveCredentials(string) error
}

type CredentialsStorageBackup interface {
	BackupStorage() (io.Reader, error)
}

type CredentialsStorageRestore interface {
	RestoreStorage(io.Reader) error
}

type MutableKey interface {
	UpdateKey([]byte) error
}
