# network-controller
A network controller implementation over websockets. This is the client/server implementation of a layer-2 software switch able to route packets over websockets connections. The network-controller is built on top of Linux's tun/tap device.

* [network-controller](#network-controller)
   * [Build and Install](#build-and-install)
   * [Adm tool](#adm-tool)
   * [Pre-work](#pre-work)
      * [Generate tls certificate](#generate-tls-certificate)
      * [Get token](#get-token)
   * [Configuration](#configuration)
      * [Download](#download)
      * [Network forwarding](#network-forwarding)
   * [network-controller-explain](#network-controller-explain)

## Build and Install

[build doc](./build/build.md)

## Adm tool

[ncadm](https://github.com/Litekube/ncadm), a commond-line tool to control node join to litekube network-controller

## Pre-work

### Generate tls certificate

[certs generation script](./build/gen_certs.sh)

```shell
cd ./build
# tls certificate dir
# network: ./certs/init/test1    grpc: ./certs/init/test2
# $ip(demo:101.43.253.110) is the host public ip or addressable private ip
sh gen_certs.sh $ip

# modify ./cmd/network-controller/server.yml
networkCertDir: /root/go_project/network-controller/certs/init/test1/
grpcCertDir: /root/go_project/network-controller/certs/init/test2/
```

### Get token

```shell
# ./cmd/ncadm compile ncadm
$ go build -o ncadm .

# generate no-expire bootstrap-token
$ ./ncadm create-bootstrap-token --life=-1

------------------------------------------------
network-controller:
    token: 2283a030cbd54b90@101.43.253.110:6439
    ExpireMsg: no expire
------------------------------------------------

# get node-token & network+grpc clients certs
# --network-certs-dir/--grpc-certs-dir is the directory where client certs store
$ ./ncadm get-token --bootstrap-token=2283a030cbd54b90 --network-certs-dir=/root/go_project/network-controller/certs/init/gen/network --grpc-certs-dir=/root/go_project/network-controller/certs/init/gen/grpc

------------------------------------------------
network-controller:
    BootstrapToken: 2283a030cbd54b90
    NodeToken: 5f5e4ced3bd44ca1
    NetworkServerIp: 101.43.253.110
    NetworkServerPort: 6441
    GrpcServerIp: 10.1.1.1
    GrpcServerPort: 6440
    NetworkCertsDir: /root/go_project/network-controller/certs/init/gen/network
    GrpcCertsDir: /root/go_project/network-controller/certs/init/gen/grpc
------------------------------------------------
```



```shell
# modify ./cmd/network-controller/client.yml
networkCertDir: /root/go_project/network-controller/certs/init/gen/network/
token: 5f5e4ced3bd44ca1
```

## Configuration & Run

There are two config files to distinguish between [client](./cmd/network-controller/client.yml) and [server](./cmd/network-controller/server.yml).

To start server/client, execute the following command:

```shell
cd ./cmd/network-controller
# server
network-controller --config server.yml
# client
network-controller --config client.yml
```

### Download

[release](https://github.com/Litekube/network-controller/releases)

### Network forwarding
On the server the IP forwarding is needed. First we need to be sure that IP forwarding is enabled.
Very often this is disabled by default. This is done by running the following command line as root：

```shell
sysctl -w net.ipv4.ip_forward=1
iptables -t nat -A POSTROUTING -j MASQUERADE
```

So, lets look at the iptables rules required for this to work.
```shell
# Allow TUN interface connections to network server
iptables -A INPUT -i tun0 -j ACCEPT

# Allow TUN interface connections to be forwarded through other interfaces
iptables -A FORWARD -i tun0 -j ACCEPT

iptables -t nat -A POSTROUTING -o tun0 -j MASQUERADE

iptables -A FORWARD -i eth0 -o tun0 -j ACCEPT

iptables -A FORWARD -i tun0 -o eth0 -m state --state RELATED,ESTABLISHED -j ACCEPT
```

## network-controller-explain doc

if you want to know more about this project, please look at :

- [Network-controller Network Part Explaination doc](docs/network-controller-explain.md)
- [API Reference doc](docs/API-explain.md)
- [PRD & System Design doc](docs/design-explain.md)
- [Usage Demo & doc](docs/demo-usage.md)

