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
package vpn

import (
	"litekube-vpn/utils"
)

var logger = utils.GetLogger()

var MTU = 1400

const NodeTokenKey = "node-token"

const (
	IFACE_BUFSIZE = 2048
)

const (
	STATE_IDLE      = -1
	STATE_INIT      = 1
	STATE_CONNECT   = 2
	STATE_CONNECTED = 3
)

const (
	STATUS_OK         = "200"
	STATUS_BADREQUEST = "400"
	STATUS_ERR        = "500"
)

const (
	MESSAGE_OK = "ok"
)
