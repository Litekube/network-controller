简体中文 | [English](./network-controller-explain.md)

# network-controller-explain

* [network-controller-explain](#network-controller-explain)
   * [network C/S交互关系](#network-cs交互关系)
      * [C/S通信](#cs通信)
      * [虚拟网卡 &amp; 路由](#虚拟网卡--路由)
   * [Code detail](#code-detail)
      * [Interface.go](#interfacego)

## network C/S交互关系

### C/S通信

- 发送请求：从网卡读message数据，通过data channel，发送websocket请求
- 接收请求：c/s状态为STATE_CONNECTED后，收到websocket请求，通过toIface channel，向网卡写数据

> C/S状态

- Client/Server初次连接，通过一次connect state的message后，双方的状态从init变为connected
- 往后可以正常通信

![aim3-状态](https://tva1.sinaimg.cn/large/e6c9d24ely1h163y8cmbdj219e0hmta6.jpg)

- server connection/clien state：init or connected
- message state
  - connect：client发给server初次建立连接请求
  - connected：已成功建立连接

### 虚拟网卡 & 路由

- server端
  - 目的地址10.1.1.0/24的请求 gateway到10.1.1.2，在gateway到0.0.0.0 flag=host 出去

```shell
# ifconfig
tun0
inet 10.1.1.1  netmask 255.255.255.255  destination 10.1.1.2
# route
Destination     Gateway         Genmask         Flags Metric Ref    Use Iface
10.1.1.0        10.1.1.2        255.255.255.0   UG    0      0        0 tun0
10.1.1.2        0.0.0.0         255.255.255.255 UH    0      0        0 tun0
```

- client端
  - 目的地址10.1.1.0/24的请求 gateway到10.1.1.4，在gateway到0.0.0.0 flag=host 出去
  - 目的地址101.43.253.110(network-server)的请求，gateway到{gateway}（即内网地址 ens160网卡），通过ens160出去

```shell
# ifconfig
tun0
inet 10.1.1.3  netmask 255.255.255.255  destination 10.1.1.4
# route
Destination     Gateway         Genmask         Flags Metric Ref    Use Iface
10.1.1.0        10.1.1.4        255.255.255.0   UG    0      0        0 tun0
10.1.1.4        0.0.0.0         255.255.255.255 UH    0      0        0 tun0
101.43.253.110  gateway         255.255.255.255 UGH   0      0        0 ens160
```

## Code detail

### Interface.go

- newTun

> 创建tun0 interface

```
iface, err = water.New(water.Config{})
```

> ip是iproute2软件包里面的一个强大的网络配置工具，它能够替代一些传统的网络管理工具（route、ifconfig等

```shell
# ifconfig tun0 up
# 改变网络设备MTU(最大传输单元)的值 mtu=1400
# 改变设备传输队列的长度 qlen=100
ip link set dev tun0 up mtu 1400 qlen 100
```

- setTunIP

```shell
# 配置网卡tun0 10.1.1.1/32
# 使用点对点连接时对端的协议地址 peer 10.1.1.2
ip addr add dev tun0 local 10.1.1.1 peer 10.1.1.2

# result 
# ifconfig tun0
inet 10.1.1.1  netmask 255.255.255.255  destination 10.1.1.2
# route 目的10.1.1.2 host out
Destination     Gateway         Genmask         Flags Metric Ref    Use Iface
10.1.1.2        0.0.0.0         255.255.255.255 UH    0      0        0 tun0
```

```shell
# add route
ip route add 10.1.1.0/24 via 10.1.1.2 dev tun0

# result 目的10.1.1.0/24，gateway路由到10.1.1.2
Destination     Gateway         Genmask         Flags Metric Ref    Use Iface
10.1.1.0        10.1.1.2        255.255.255.0   UG    0      0        0 tun0
```

- client addRoute

```shell
# 添加101.43.253.110  gateway         255.255.255.255 UGH   0      0        0 ens160
ip -4 r a 101.43.253.110/32 via 192.168.107.2 dev ens160
```
