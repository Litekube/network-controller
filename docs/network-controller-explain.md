English | [简体中文](./network-controller-explain_CN.md)

# network-controller-explain

* [network-controller-explain](#network-controller-explain)
   * [network C/S interaction relationship](#network-cs-interaction-relationship)
      * [C/S communication](#cs-communication)
      * [Virtual NICs &amp; Routing](#virtual-nics--routing)
   * [Code detail](#code-detail)
      * [Interface.go](#interfacego)

## network C/S interaction relationship

### C/S communication

- Sending a request: Reads message data from the network adapter and sends a WebSocket request through the Data Channel
- Receiving requests: After the C/S status is STATE_CONNECTED, a WebSocket request is received and data is written to a network adapter through the toIface Channel

> C/S state

- After a connect state message is sent to the Client/Server for the first time, the state changes from STATE_INIT to STATE_CONNECTED
- Communication will be normal from now on

![aim3-状态](https://tva1.sinaimg.cn/large/e6c9d24ely1h163y8cmbdj219e0hmta6.jpg)

- server connection/clien state：STATE_INIT or STATE_CONNECTED
- message state
  - STATE_CONNECT：The client sends an initial connection request to the server
  - STATE_CONNECTED：The connection has been established

### Virtual NICs & Routing

- server
  - The request destined for 10.1.1.0/24 is routed to gateway 10.1.1.2 and then 0.0.0.0(flag=host)

```shell
# ifconfig
tun0
inet 10.1.1.1  netmask 255.255.255.255  destination 10.1.1.2
# route
Destination     Gateway         Genmask         Flags Metric Ref    Use Iface
10.1.1.0        10.1.1.2        255.255.255.0   UG    0      0        0 tun0
10.1.1.2        0.0.0.0         255.255.255.255 UH    0      0        0 tun0
```

- client
  - The request destined for 10.1.1.0/24 is routed to gateway 10.1.1.4 and then 0.0.0.0(flag=host)
  - The request of the destination address 101.43.253.110 (network-server), routed to {gateway} (that is, the intranet address ens160 nic), go out through ens160

```shell
# ifconfig
tun0
inet 10.1.1.3  netmask 255.255.255.255  destination 10.1.1.4

# route 
# 101.43.253.110 is the public address
Destination     Gateway         Genmask         Flags Metric Ref    Use Iface
10.1.1.0        10.1.1.4        255.255.255.0   UG    0      0        0 tun0
10.1.1.4        0.0.0.0         255.255.255.255 UH    0      0        0 tun0
101.43.253.110  gateway         255.255.255.255 UGH   0      0        0 ens160
```

## Code detail

### Interface.go

- newTun

> create tun0 interface

```
iface, err = water.New(water.Config{})
```

> IP is a powerful network configuration tool in the iproute2 software package, which can replace some traditional network management tools (route, ifconfig, etc.).

```shell
# ifconfig tun0 up
# Change the MTU(Maximum transmission unit) value of the network device MTU to 1400
# Change the length of the transmission queue qlen=100
ip link set dev tun0 up mtu 1400 qlen 100
```

- setTunIP

```shell
# interface tun0 10.1.1.1/32
ip addr add dev tun0 local 10.1.1.1 peer 10.1.1.2

# result 
# ifconfig tun0
inet 10.1.1.1  netmask 255.255.255.255  destination 10.1.1.2

# route Destination 10.1.1.2 host out
Destination     Gateway         Genmask         Flags Metric Ref    Use Iface
10.1.1.2        0.0.0.0         255.255.255.255 UH    0      0        0 tun0
```

```shell
# add route
ip route add 10.1.1.0/24 via 10.1.1.2 dev tun0

# result Destination 10.1.1.0/24, gateway to 10.1.1.2
Destination     Gateway         Genmask         Flags Metric Ref    Use Iface
10.1.1.0        10.1.1.2        255.255.255.0   UG    0      0        0 tun0
```

- client addRoute

```shell
# add 101.43.253.110  gateway         255.255.255.255 UGH   0      0        0 ens160
ip -4 r a 101.43.253.110/32 via 192.168.107.2 dev ens160
```
