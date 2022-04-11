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
 * Author: Lukasz Zajaczkowski <zreigz@gmail.com>
 *
 */

package vpn

import (
	"net"

	"fmt"
	"github.com/songgao/water"
	//. "github.com/zreigz/ws-vpn/vpn/utils"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	. "ws-vpn/vpn/utils"

	"golang.org/x/net/ipv4"
)

type VpnServer struct {
	// config
	cfg ServerConfig
	// interface
	iface *water.Interface
	// subnet
	ipnet *net.IPNet
	// IP Pool
	ippool *VpnIpPool
	// client peers, key is the mac address, value is a HopPeer record
	// Registered clients clientip-connection
	clients map[string]*connection
	// Register requests
	register chan *connection
	// Unregister requests
	unregister chan *connection
	outData    *Data
	inData     chan *Data
	toIface    chan []byte
}

func NewServer(cfg ServerConfig) error {
	var err error

	if cfg.MTU != 0 {
		MTU = cfg.MTU
	}

	vpnServer := &VpnServer{}
	vpnServer.cfg = cfg
	vpnServer.ippool = &VpnIpPool{}

	iface, err := newTun("")
	if err != nil {
		return err
	}
	vpnServer.iface = iface

	// vpnaddr = 10.1.1.1/24
	ip, subnet, err := net.ParseCIDR(cfg.VpnAddr)
	err = setTunIP(iface, ip, subnet)
	if err != nil {
		return err
	}
	vpnServer.ipnet = &net.IPNet{ip, subnet.Mask}
	vpnServer.ippool.subnet = subnet

	go vpnServer.cleanUp()

	go vpnServer.run()

	vpnServer.register = make(chan *connection)
	vpnServer.unregister = make(chan *connection)
	vpnServer.clients = make(map[string]*connection)
	// no use
	vpnServer.inData = make(chan *Data, 100)
	vpnServer.toIface = make(chan []byte, 100)

	vpnServer.handleInterface()

	// http handle for client to connect
	http.HandleFunc("/ws", vpnServer.serveWs)
	addr := fmt.Sprintf(":%d", vpnServer.cfg.Port)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		logger.Panicf("ListenAndServe: %+v" + err.Error())
	}
	logger.Infof("ListenAndServe: %+v", addr)

	return nil

}

func (server *VpnServer) serveWs(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	// client http to ws
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error(err)
		return
	}
	NewConnection(ws, server)
}

func (server *VpnServer) run() {
	for {
		select {
		case c := <-server.register:
			// add to clients
			logger.Infof("Connection registered: %+v", c.ipAddress.IP.String())
			server.clients[c.ipAddress.IP.String()] = c
			break

		case c := <-server.unregister:
			// remove from clients
			// close connection data channel
			// release client ip
			clientIP := c.ipAddress.IP.String()
			_, ok := server.clients[clientIP]
			if ok {
				delete(server.clients, clientIP)
				close(c.data)
				if c.ipAddress != nil {
					server.ippool.release(c.ipAddress.IP)
				}
				logger.Infof("unregister Connection: %+v, current active clients number: %+v", c.ipAddress.IP, len(server.clients))
			}
			break
		}
	}
}

func (server *VpnServer) handleInterface() {
	// network packet to interface
	go func() {
		for {
			hp := <-server.toIface
			logger.Debug("Write to interface")
			_, err := server.iface.Write(hp)
			if err != nil {
				logger.Error(err.Error())
				return
			}

		}
	}()

	// interface to network packet
	go func() {
		packet := make([]byte, IFACE_BUFSIZE)
		for {
			plen, err := server.iface.Read(packet)
			if err != nil {
				logger.Error(err)
				break
			}
			header, _ := ipv4.ParseHeader(packet[:plen])
			logger.Debugf("Try sending: %+v", header)
			clientIP := header.Dst.String()
			client, ok := server.clients[clientIP]
			if ok {
				// config file "interconnection=false" not allowed connection between clients
				if !server.cfg.Interconnection {
					if server.isConnectionBetweenClients(header) {
						logger.Infof("Drop connection betwenn %+v and %+v", header.Src, header.Dst)
						continue
					}
				}

				logger.Debugf("Sending to client: %+v", client.ipAddress)
				client.data <- &Data{
					ConnectionState: STATE_CONNECTED,
					Payload:         packet[:plen],
				}

			} else {
				logger.Warningf("Client not found: %+v", clientIP)
			}
		}
	}()
}

func (server *VpnServer) isConnectionBetweenClients(header *ipv4.Header) bool {

	// srcip!= server ip & desip=one client ip
	if header.Src.String() != header.Dst.String() && header.Src.String() != server.ipnet.IP.String() && server.ippool.subnet.Contains(header.Dst) {
		return true
	}
	return false
}

// server exit gracefully
func (server *VpnServer) cleanUp() {

	c := make(chan os.Signal, 1)
	// watch ctrl+c or kill pid
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	logger.Debug("clean up")
	// close all client connection
	for key, client := range server.clients {
		client.ws.Close()
		delete(server.clients, key)
	}

	// code zero indicates success
	os.Exit(0)
}
