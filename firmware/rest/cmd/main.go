package main

import (
	"flag"
	"fmt"

	"github.com/jdevelop/passkeeper/firmware"
	"github.com/jdevelop/passkeeper/firmware/rest"
	"github.com/jdevelop/passkeeper/firmware/storage"
)

type storageCombined interface {
	storage.CredentialsStorageList
	storage.CredentialsStorageRead
	storage.CredentialsStorageRemove
	storage.CredentialsStorageWrite
}

type inMemory struct {
	creds []*firmware.Credentials
}

func (im *inMemory) ReadCredentials(id string) (*firmware.Credentials, error) {
	for _, v := range im.creds {
		if v.Id == id {
			return v, nil
		}
	}
	return nil, fmt.Errorf("Credential not found %s", id)
}

func (im *inMemory) WriteCredentials(creds firmware.Credentials) error {
	for i, v := range im.creds {
		if v.Id == creds.Id {
			im.creds[i] = &creds
			return nil
		}
	}
	im.creds = append(im.creds, &creds)
	return nil
}

func (im *inMemory) ListCredentials() ([]firmware.Credentials, error) {
	res := make([]firmware.Credentials, 0)
	for _, v := range im.creds {
		res = append(res, *v)
	}
	return res, nil
}

func (im *inMemory) RemoveCredentials(id string) error {
	for i, v := range im.creds {
		if v.Id == id {
			last := len(im.creds) - 1
			im.creds[i], im.creds[last] = im.creds[last], im.creds[i]
			im.creds = im.creds[:last]
			return nil
		}
	}
	return nil
}

var (
	host = flag.String("host", "localhost", "host to listen on")
	port = flag.Int("port", 8081, "port to listen on")
)

func main() {
	flag.Parse()

	initCreds := make([]*firmware.Credentials, 2)
	for i := range initCreds {
		initCreds[i] = &firmware.Credentials{
			Id:      fmt.Sprintf("%d", i),
			Service: fmt.Sprintf("service %d", i),
			Secret:  fmt.Sprintf("secret %d", i),
			Comment: fmt.Sprintf("commend %d", i),
		}
	}

	storage := inMemory{
		creds: initCreds,
	}

	fmt.Printf("Starting REST service at http://%s:%d\n", *host, *port)

	rest.Start(*host, *port, &storage, func() {})
}
