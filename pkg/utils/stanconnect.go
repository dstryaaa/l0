package utils

import (
	"log"

	"github.com/nats-io/stan.go"
)

var sc stan.Conn
var err error

func ConnectToStan() (stan.Conn, error) {
	sc, err = stan.Connect("test-cluster", "unique-client-id")
	if err != nil {
		log.Fatal(err)
	}
	return sc, err
}
