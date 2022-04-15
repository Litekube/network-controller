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
package config

import (
	"errors"
	"github.com/scalingdata/gcfg"
)

// server.ini / client.ini

// Server Config
type ServerConfig struct {
	Port            int
	GrpcPort        int
	ListenAddr      string
	VpnAddr         string
	MTU             int
	Interconnection bool
}

// Client Config
type ClientConfig struct {
	Server          string
	Port            int
	MTU             int
	Token           string
	RedirectGateway bool
}

type VpnConfig struct {
	Default struct {
		Mode string
	}
	Server ServerConfig
	Client ClientConfig
}

// return server/client config
func ParseConfig(filename string) (interface{}, error) {
	cfg := &VpnConfig{}
	err := gcfg.ReadFileInto(cfg, filename)
	if err != nil {
		return nil, err
	}
	switch cfg.Default.Mode {
	case "server":
		return cfg.Server, nil
	case "client":
		return cfg.Client, nil
	default:
		return nil, errors.New("Wrong config data")
	}
}
