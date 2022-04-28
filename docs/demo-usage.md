# network-controller使用demo

* [network-controller使用demo](#network-controller使用demo)
    * [network部分](#network部分)
    * [grpc部分](#grpc部分)
    * [一次完成执行流程](#一次完成执行流程)

## network部分

- server/client执行详见network/network.go
    - 根据yml配置文件的不同，可以自动鉴别server or client并启动
    - 配置文件字段含义见server.yml/client.yml注册

```shell
# network server & grpc server启动
go run network/network.go server.yml

# network client启动
go run network/network.go client.yml
```

## grpc部分

- 详见grpc/grpc_client.go
    - 运行前应该先执行 `go run network/network.go server.yml`

```shell
cd grpc
# get bootstrap token
go test -v -run TestGetBootstrapToken

# 注册 得到token=a9f683a2d05b4957，以及返回的grpc+network证书
go test -v -run TestGetToken

# 查询连接状态
go test -v -run TestCheckConnState 

# 取消注册（取消绑定 & 断连
go test -v -run TestUnRegister 
```

## 一次完整执行流程

```shell
# 启动network & grpc server
go run network/network.go server.yml

cd grpc
# get bootstrap token
go test -v -run TestGetBootstrapToken

# 注册 得到token=a9f683a2d05b4957，以及返回的grpc+network证书
go test -v -run TestGetToken

# 根据token修改client.yml，启动client
go run network/network.go client.yml

# 查询连接状态
go test -v -run TestCheckConnState 

# 取消注册（取消绑定 & 断连
go test -v -run TestUnRegister 
```

