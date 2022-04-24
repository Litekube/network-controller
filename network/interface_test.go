package network

import (
	"log"
	"testing"
)

func TestGetNetGateway(t *testing.T) {
	// fix file path "./test.txt" before test
	gateway, dev, err := GetNetGateway()
	log.Printf(gateway)
	log.Printf(dev)
	if err != nil {
		log.Printf(err.Error())
	}
}
