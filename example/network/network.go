package main

import (
	"flag"
	"github.com/Litekube/network-controller/config"
	"github.com/Litekube/network-controller/network"
	"github.com/Litekube/network-controller/utils"
	"os"
)

var debug bool
var cfgFile string

func main() {
	flag.BoolVar(&debug, "debug", false, "Provide debug info")
	flag.StringVar(&cfgFile, "config", "", "config file")
	flag.Parse()

	utils.InitLogger()
	utils.SetLoggerLevel(debug)

	logger := utils.GetLogger()

	checkerr := func(err error) {
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
	}

	if cfgFile == "" {
		cfgFile = flag.Arg(0)
	}

	logger.Infof("using config file: %+v", cfgFile)

	icfg, err := config.ParseConfig(cfgFile)
	logger.Debug(icfg)
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
		logger.Error("Invalid config file")
	}
}
