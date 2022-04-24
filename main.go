/*
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 * Author: wanna <wananzjx@163.com>
 *
 */

package main

import (
	"flag"
	"github.com/Litekube/network-controller/config"
	"github.com/Litekube/network-controller/utils"
	client "github.com/Litekube/network-controller/vpn"
	server "github.com/Litekube/network-controller/vpn"
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

	//maxProcs := runtime.GOMAXPROCS(0)
	//if maxProcs < 2 {
	//	runtime.GOMAXPROCS(2)
	//}

	switch cfg := icfg.(type) {
	case config.ServerConfig:
		vpnServer := server.NewServer(cfg)
		err = vpnServer.Run()
		checkerr(err)
	case config.ClientConfig:
		client := client.NewClient(cfg)
		err := client.Run()
		checkerr(err)
	default:
		logger.Error("Invalid config file")
	}
}
