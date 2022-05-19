English | [简体中文](./demo-usage_CN.md)

# network-controller usage demo

* [network-controller usage demo](#network-controller-usage-demo)
   * [Pre-work](#pre-work)
   * [network part](#network-part)
   * [grpc part](#grpc-part)
   * [A complete execution process](#a-complete-execution-process)

## Pre-work

[Pre-work certs+token](https://github.com/Litekube/network-controller/tree/main#pre-work)

## network part

- For details on network server/client execution, see network/network.go
    - Depending on the yaml configuration file, the server or client can be automatically identified and started
    - For the meaning of the configuration file fields, see the comments in server.yml/client.yml

```shell
# run network server & grpc server
go run network/network.go server.yml

# run network client
go run network/network.go client.yml
```

## grpc part

- For details on grpc client execution, see grpc/grpc_client.go
    - Before running, execute: `go run network/network.go server.yml`

```shell
cd grpc
# get bootstrap token
go test -v -run TestGetBootstrapToken

# Register to get token= a9f683a2d05b4957, and return GRPC + Network certificate
go test -v -run TestGetToken

# Query connection Status
go test -v -run TestCheckConnState 

# Unregister (unbind & disconnect)
go test -v -run TestUnRegister 
```

## A complete execution process

```shell
# run network & grpc server
go run network/network.go server.yml

cd grpc
# get bootstrap token
go test -v -run TestGetBootstrapToken

# Register to get token= a9f683a2d05b4957, and return GRPC + Network certificate
go test -v -run TestGetToken

# Modify client.yml according to the token and start the client
go run network/network.go client.yml

# Query connection Status
go test -v -run TestCheckConnState 

# Unregister (unbind & disconnect)
go test -v -run TestUnRegister 
```

