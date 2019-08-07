package storage

import (
	"github.com/jdevelop/passkeeper/firmware"
)

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
