# ncadm

[ncadm](https://github.com/Litekube/ncadm), a commond-line tool to control node join to litekube network-controller

## build & installation

[build doc](../build/build.md)

## usage

```shell
$ ./ncadm -h

NAME:
   ncadm - ncadm, a commond-line tool to control node join to litekube network-controller

USAGE:
   ncadm [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
   create-bootstrap-token  create network bootstrap token info
   get-token               get grpc server ip/port/certs
   check-conn-state        check network conn state
   unregister              close network connection unregister bind ip
   check-health            check health of control and bootstrap grpc
   help, h                 Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --ip value        leader host ip (default: "127.0.0.1")
   --port value      network grpc control port (default: "6440")
   --bootport value  network grpc bootstrap control port (default: "6439")
   --cacert value    ca cert filepath of network grpc server (default: "/root/.litekube/nc/certs/grpc/ca.pem")
   --cert value      client cert filepath of network grpc server (default: "/root/.litekube/nc/certs/grpc/client.pem")
   --key value       client key filepath of network grpc server (default: "/root/.litekube/nc/certs/grpc/client-key.pem")
   --help, -h        show help (default: false)
   --version, -v     print the version (default: false)
   
   
$ ncadm-Linux-x86_64-0.1.0 -v
version.Info{NcAdm:"0.1.0", NetworkController:"0.1.0", GitVersion:"git version 1.8.3.1", GitBranch:"main", GitCommit:"7853dcb364e64e5e9806ec918c9152ed7f37fa7b", BuildDate:"2022-05-17T02:17:25Z", GoVersion:"go1.17.8", Compiler:"gc", Platform:"linux/amd64"}
```

### CreateBootstrapToken

```shell
$ ./ncadm --ip 101.43.253.110 create-bootstrap-token --life=10

------------------------------------------------
network-controller:
    token: 7ce776409acd4d66@101.43.253.110:6439
    ExpireMsg: expire in 10 min
------------------------------------------------
```

### CheckConnState

```shell
$ ./ncadm --ip 101.43.253.110 check-conn-state  --node-token=c9168a88d15b4a18

------------------------------------------------
network-controller:
    node-token: c9168a88d15b4a18
    BindIp: 10.1.1.3
    ConnState: UnConnected
------------------------------------------------
```

### GetToken

```shell
$ ./ncadm --ip 101.43.253.110 get-token --bootstrap-token=33df5c33f42f4960 --network-certs-dir=/Users/zhujianxing/nc/certs/test/network --grpc-certs-dir=/Users/zhujianxing/nc/certs/test/grpc

------------------------------------------------
network-controller:
    node-token: 33df5c33f42f4960
    NetworkServerIp: 101.43.253.110
    NetworkServerPort: 6441
    GrpcServerIp: 10.1.1.1
    GrpcServerPort: 6440
    NetworkCertsDir: /Users/zhujianxing/nc/certs/test/network
    GrpcCertsDir: /Users/zhujianxing/nc/certs/test/grpc
------------------------------------------------
```

### CheckHealth

```shell
$ ./ncadm --ip 101.43.253.110 check-health

------------------------------------------------
network-controller:
    control grpc client health: Health
    bootstrap grpc client health: Health
------------------------------------------------
```

### UnRegister

```shell
$ ./ncadm --ip 101.43.253.110 unregister --node-token=c9168a88d15b4a18

------------------------------------------------
network-controller:
    node-token: c9168a88d15b4a18
    Success: true
------------------------------------------------
```

