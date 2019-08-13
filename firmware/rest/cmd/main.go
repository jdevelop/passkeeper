package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/jdevelop/passkeeper/firmware/pass"
	"github.com/jdevelop/passkeeper/firmware/rest"
	"github.com/jdevelop/passkeeper/firmware/storage"
)

var (
	host     = flag.String("host", "localhost", "host to listen on")
	port     = flag.Int("port", 8081, "port to listen on")
	passfile = flag.String("passfile", "/tmp/encrypted.storage", "path to the passwords file")
)

func main() {
	flag.Parse()

	s, err := storage.NewPlainText(*passfile, []byte("passw0rd12345678"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Starting REST service at http://%s:%d\n", *host, *port)

	rest.Start(*host, *port, s, pass.NewPasswordGenerator(8), func() {})
}
