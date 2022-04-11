package vpn

import (
	"log"
	"testing"
)

func TestGetNetGateway(t *testing.T) {
	// fix file path before test
	gateway, dev, err := GetNetGateway()
	log.Printf(gateway)
	log.Printf(dev)
	if err != nil {
		log.Printf(err.Error())
	}
}
