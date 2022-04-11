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
	"errors"
	"net"
	"sync/atomic"
)

/*
assign unique ip for client
todo persist data to db
*/

type VpnIpPool struct {
	subnet *net.IPNet
	pool   [127]int32 // map
}

var poolFull = errors.New("IP Pool Full")

// get an empty ip
func (p *VpnIpPool) next() (*net.IPNet, error) {
	found := false
	var i int
	// server take x.1 & x.2, begin from 3
	for i = 3; i < 255; i += 2 {
		// CAS sync
		if atomic.CompareAndSwapInt32(&p.pool[i], 0, 1) {
			found = true
			break
		}
	}
	if !found {
		return nil, poolFull
	}

	// assign ip+mask
	ipnet := &net.IPNet{
		make([]byte, 4),
		make([]byte, 4),
	}
	copy([]byte(ipnet.IP), []byte(p.subnet.IP))
	copy([]byte(ipnet.Mask), []byte(p.subnet.Mask))
	ipnet.IP[3] = byte(i) // found=true
	return ipnet, nil
}

// release ip
func (p *VpnIpPool) release(ip net.IP) {
	defer func() {
		// recover only work in defer part
		// if normal, return nil
		// if panic, return panic err and recover normal,continue to execute
		if err := recover(); err != nil {
			logger.Errorf("release err:%v", err)
		}
	}()

	logger.Infof("releasing ip: %+v", ip)
	i := ip[3]
	p.pool[i] = 0
}
