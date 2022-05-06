package main

import (
	"flag"
	"fmt"
	"github.com/Litekube/network-controller/config"
	"github.com/Litekube/network-controller/network"
	"os"
)

var debug bool
var cfgFile string

func main() {
	flag.BoolVar(&debug, "debug", false, "Provide debug info")
	flag.StringVar(&cfgFile, "config", "", "config file")
	flag.Parse()

	checkerr := func(err error) {
		if err != nil {
			os.Exit(1)
		}
	}

	if cfgFile == "" {
		cfgFile = flag.Arg(0)
	}

	icfg, err := config.ParseConfig(cfgFile)
	checkerr(err)

	switch cfg := icfg.(type) {
	case config.ServerConfig:
		server := network.NewServer(cfg)
		err = server.Run()
		checkerr(err)
	case config.ClientConfig:
		client := network.NewClient(cfg)
		err := client.Run()
		checkerr(err)
	default:
	}
	fmt.Println("main exit")
}
