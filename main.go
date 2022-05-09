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
	"fmt"
	"github.com/Litekube/network-controller/config"
	client "github.com/Litekube/network-controller/network"
	server "github.com/Litekube/network-controller/network"
	"os"
)

var debug bool
var cfgFile string

func main() {
	flag.BoolVar(&debug, "debug", false, "Provide debug info")
	flag.StringVar(&cfgFile, "config", "", "config file")
	flag.Parse()

	if cfgFile == "" {
		cfgFile = flag.Arg(0)
	}

	icfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		os.Exit(1)
	}

	//maxProcs := runtime.GOMAXPROCS(0)
	//if maxProcs < 2 {
	//	runtime.GOMAXPROCS(2)
	//}

	switch cfg := icfg.(type) {
	case config.ServerConfig:
		networkServer := server.NewServer(cfg)
		//err = networkServer.Run()
		//if err != nil {
		//	logger.Error(err.Error())
		//	os.Exit(1)
		//}
		go func() {
			err = networkServer.Run()
			fmt.Println(err)
		}()

		defer networkServer.Stop()
	case config.ClientConfig:
		networkClient := client.NewClient(cfg)
		//err := networkClient.Run()
		//if err != nil {
		//	logger.Error(err.Error())
		//	os.Exit(1)
		//}
		go func() {
			err = networkClient.Run()
			fmt.Println(err)
		}()
		defer networkClient.Wait()
	default:
	}
	fmt.Println("main exit")
}
