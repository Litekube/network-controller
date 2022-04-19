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
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

// server.yml / client.yml

// Server Config
type ServerConfig struct {
	Port        int    `yaml:"port"`
	VpnCertDir  string `yaml:"vpnCertDir"`
	GrpcPort    int    `yaml:"grpcPort"`
	GrpcCertDir string `yaml:"grpcCertDir"`

	ListenAddr      string `yaml:"listenAddr"`
	VpnAddr         string `yaml:"vpnAddr"`
	MTU             int    `yaml:"mtu"`
	Interconnection bool   `yaml:"interconnection"`
}

// Client Config
type ClientConfig struct {
	VpnCertDir      string `yaml:"vpnCertDir"`
	ServerAddr      string `yaml:"serverAddr"`
	Port            int    `yaml:"port"`
	MTU             int    `yaml:"mut"`
	Token           string `yaml:"token"`
	RedirectGateway bool   `yaml:"redirectGateway"`
}

type VpnConfig struct {
	Mode   string       `yaml:"mode"`
	Server ServerConfig `yaml:"server"`
	Client ClientConfig `yaml:"client"`
}

// return server/client config
func ParseConfig(filename string) (interface{}, error) {
	cfg := &VpnConfig{}

	File, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("fail to read file: %v", err)
	}
	err = yaml.Unmarshal(File, &cfg)
	if err != nil {
		return nil, err
	}

	switch cfg.Mode {
	case "server":
		return cfg.Server, nil
	case "client":
		return cfg.Client, nil
	default:
		return nil, errors.New("Wrong config data")
	}
}
