## gRPC接口说明

* [gRPC接口说明](#grpc接口说明)
  * [概述](#概述)
  * [grpcurl工具](#grpcurl工具)
  * [接口列表](#接口列表)
       * [Bootstrap注册 GetBootStrapToken](#bootstrap注册-getbootstraptoken)
       * [节点注册 GetToken](#节点注册-gettoken)
       * [节点取消注册 UnRegister](#节点取消注册-unregister)
       * [检查连接状态 CheckConnState](#检查连接状态-checkconnstate)
       * [获取连接节点ip GetRegistedIp](#获取连接节点ip-getregistedip)

### 概述

基于tcp gRPC+protobuf实现通信交互服务，支持tls安全通信

```
service LiteKubeNCService {
  rpc GetBootStrapToken(GetBootStrapTokenRequest) returns (GetBootStrapTokenResponse) {}
  rpc GetToken(GetTokenRequest) returns (GetTokenResponse) {}
  rpc CheckConnState(CheckConnStateRequest) returns (CheckConnResponse){}
  rpc UnRegister(UnRegisterRequest) returns (UnRegisterResponse){}
  rpc GetRegistedIp(GetRegistedIpRequest) returns (GetRegistedIpResponse){}
}
```

- 状态码code说明

| code枚举类型      | 值   | 含义               |
| ----------------- | ---- | ------------------ |
| STATUS_OK         | 200  | 成功               |
| STATUS_BADREQUEST | 400  | 客户端参数不规范   |
| STATUS_ERR        | 500  | 服务器内部逻辑错误 |

### grpcurl工具

> installation

- mac

```shell
brew install grpcurl
```

- linux

```shell
wget https://github.com/fullstorydev/grpcurl/releases/download/v1.8.5/grpcurl_1.8.5_linux_x86_64.tar.gz 
tar -xvf grpcurl_1.8.5_linux_x86_64.tar.gz 
chmod +x grpcurl 
cp grpcurl /usr/sbin/
```

> usage

- without tls

```shell
# 查看服务列表
grpcurl -plaintext 101.43.253.110:6440 list pb.LiteKubeNCService

# 查看请求
grpcurl -plaintext 101.43.253.110:6440 describe pb.LiteKubeNCService.HelloWorld

# 查看参数
grpcurl -plaintext 101.43.253.110:6440 describe pb.HelloWorldRequest 
grpcurl -plaintext 101.43.253.110:6440 describe pb.HelloWorldResponse

# grpc调用
grpcurl -d '{"token": "b52f93d3f0ec4be7"}' -plaintext 101.43.253.110:6440 pb.LiteKubeNCService.CheckConnState

grpcurl -d '{"token": "b52f93d3f0ec4be7"}' -plaintext 101.43.253.110:6440 pb.LiteKubeNCService.UnRegister

grpcurl -d '{"token": "b52f93d3f0ec4be7"}' -plaintext 101.43.253.110:6440 pb.LiteKubeNCService.GetRegistedIp

grpcurl -d '{"bootStrapToken": "deac5f329feb4729"}' -plaintext 101.43.253.110:6439 pb.LiteKubeNCBootstrapService.GetToken

grpcurl -d '{"expireTime": 10}' -plaintext 101.43.253.110:6440 pb.LiteKubeNCService.GetBootStrapToken

```

- support tls

```shell
# grpc调用
grpcurl -d '{"token": "b52f93d3f0ec4be7"}' -cacert ca.pem -cert client.pem -key client-key.pem 101.43.253.110:6440 pb.LiteKubeNCService.CheckConnState

grpcurl -d '{"token": "b52f93d3f0ec4be7"}' -cacert ca.pem -cert client.pem -key client-key.pem 101.43.253.110:6440 pb.LiteKubeNCService.UnRegister

grpcurl -d '{"token": "b52f93d3f0ec4be7"}' -cacert ca.pem -cert client.pem -key client-key.pem 101.43.253.110:6440 pb.LiteKubeNCService.GetRegistedIp

grpcurl -d '{"bootStrapToken": "deac5f329feb4729"}' -insecure 101.43.253.110:6439 pb.LiteKubeNCBootstrapService.GetToken

grpcurl -d '{"bootStrapToken": "deac5f329feb4729"}' -cacert ca.pem -cert client.pem -key client-key.pem 101.43.253.110:6439 pb.LiteKubeNCBootstrapService.GetToken

grpcurl -d '{"expireTime": 10}' -cacert ca.pem -cert client.pem -key client-key.pem  101.43.253.110:6440 pb.LiteKubeNCService.GetBootStrapToken
```

### 接口列表

#### Bootstrap注册 GetBootStrapToken

- demo：获取一个过期时间为10min的Bootstrap token

```shell
grpcurl -d '{"expireTime": 10}' -cacert ca.pem -cert client.pem -key client-key.pem 101.43.253.110:6440 pb.LiteKubeNCService.GetBootStrapToken
```

- GetBootStrapTokenRequest参数

| **参数**   | 类型  | 含义     | 是否必须 | demo |
| ---------- | ----- | -------- | -------- | ---- |
| expireTime | int32 | 过期时间 | 否       | 10   |

- GetBootStrapTokenResponse返回数据

```json
{
  "code": "200",
  "message": "ok",
  "bootStrapToken": "deac5f329feb4729",
  "cloudIp": "101.43.253.110",
  "port": "6440"
}
```

#### 节点注册 GetToken

- demo：节点注册node-token

```shell
grpcurl -d '{"bootStrapToken": "deac5f329feb4729"}' -cacert ca.pem -cert client.pem -key client-key.pem 101.43.253.110:6440 pb.LiteKubeNCService.GetToken
```

- GetTokenRequest参数

| **参数**       | 类型   | 含义            | 是否必须 | demo             |
| -------------- | ------ | --------------- | -------- | ---------------- |
| bootStrapToken | string | bootstrap token | 是       | deac5f329feb4729 |

- GetTokenResponse返回数据（证书字段均为base64编码）

```JSON
{
  "code": "200",
  "message": "ok",
  "token": "b52f93d3f0ec4be7",
  "grpcCaCert": "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURxRENDQXBDZ0F3SUJBZ0lVYW1GNkF6L0ZCSmlrOHVRQUFuV1dvdWRFOWVJd0RRWUpLb1pJaHZjTkFRRUwKQlFBd2JERUxNQWtHQTFVRUJoTUNRMDR4RURBT0JnTlZCQWdUQjBKbGFXcHBibWN4RURBT0JnTlZCQWNUQjBKbAphV3BwYm1jeEVUQVBCZ05WQkFvVENHeHBkR1ZyZFdKbE1Rd3dDZ1lEVlFRTEV3TjJjRzR4R0RBV0JnTlZCQU1NCkQyeHBkR1ZyZFdKbFgzWndibDlqWVRBZUZ3MHlNakEwTVRZd05qVTFNREJhRncweU56QTBNVFV3TmpVMU1EQmEKTUd3eEN6QUpCZ05WQkFZVEFrTk9NUkF3RGdZRFZRUUlFd2RDWldscWFXNW5NUkF3RGdZRFZRUUhFd2RDWldscQphVzVuTVJFd0R3WURWUVFLRXdoc2FYUmxhM1ZpWlRFTU1Bb0dBMVVFQ3hNRGRuQnVNUmd3RmdZRFZRUUREQTlzCmFYUmxhM1ZpWlY5MmNHNWZZMkV3Z2dFaU1BMEdDU3FHU0liM0RRRUJBUVVBQTRJQkR3QXdnZ0VLQW9JQkFRRE8KakRHQ20vRUdLTFJoZ0tsMU8xeVEwamEvYlhjN1VpQnBoNzh6Sm92aStKanZYbVhTdGJpTElxd1o1YkdiMEdKWApLUW13bmRDVmgycUpFUXZCZ2FGc0pMdTFRY09uMU5TZ3plV3ZtYk9yNHhJQ0x1Y2QyRVZWWForNXplNmxSSjcxCjh3aHRpYjVoSlZqdTVWY1VVY2dJNStOWm5zSWdRQ2htaDcvM2RXRkh0MmR3TGxyaVpFcVVXTldLemRPeDI0MFYKeXI4QmY5UXg4MFRhYmFDVHZTUmx0TkE0TFdTYXRMTFZlV2JnMkRDMmQvcUwyMEhaQ1pPek9kT1FDc29xL3d5TwpURkZQL2FHUWFUU0lYRWZLWWtwR01oV1V2RlNoY1pObzBkcEg4czU4UkZweU5wUzJCNXcwV0NyRHRHTlVMN2tQCm1IWlBOVC9GVHhLVDNPcStJdTkvQWdNQkFBR2pRakJBTUE0R0ExVWREd0VCL3dRRUF3SUJCakFQQmdOVkhSTUIKQWY4RUJUQURBUUgvTUIwR0ExVWREZ1FXQkJUeWRsMmx2Z2tiU2cyZ3Z2djdsWlEzUGVVSkJUQU5CZ2txaGtpRwo5dzBCQVFzRkFBT0NBUUVBcnhIZk15YnpLV05jcDQ2MWppOGVWbTNqWXVldUg5ZmRXTUc2ZUxPYTFtNUZJdk5ZCm5LWUVCS1pUTEV4T1N0SklHbUZqbGpGZjNPaEhaSllmcDZnUGFFNGVuUEdDQ1Q1Wkhlcnd3dnNOL2VIVXcwR2wKWlZTSTR1OGozN3ZLZVd4aWZTaHFvc09lWHNhdW5sanlwSVQrWDBQSEx0SnRpMS9pbHlRcmVZaWNTNkIwMGhoKwpCNW1HNUpKRCs5WnVnL2hPbkQ5WXI4RU5mdnNKQzcxemdVVUxGV1NiVi96T1ZhNHlJdzhRU0pnaHFRSmR5MkpRClhGSXp6aXhseGU3NWR1UUVaNGlIWUZhblVHNVlPMFFsZ2FPazZocUp2aU1uUEhaTWorSEFPT2RPTEEwQk15RkcKWnJDMit6b2FMSERrTlExQTdQRHFoMk1VMUFOcGRqU0EvQWhYN0E9PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==",
  "grpcClientKey": "LS0tLS1CRUdJTiBFQyBQUklWQVRFIEtFWS0tLS0tCk1IY0NBUUVFSUJpWmdNd2tqNWtUS3U1US9UV3owV1ZDTVRHemtSbG1TNGl2QUpRelNOM2FvQW9HQ0NxR1NNNDkKQXdFSG9VUURRZ0FFM1RMSHphMGJNRHRvNWhhNFROcWplbi91cVVQa05GNlNGS0dkZ2J5eW1lYTNBVGhnek0wQQpWMU5ldEo3RWVmZ1dOdnZLcVdhOXlVNkMxMHBTSC9MS3JBPT0KLS0tLS1FTkQgRUMgUFJJVkFURSBLRVktLS0tLQo=",
  "grpcClientCert": "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNuVENDQVlXZ0F3SUJBZ0lJT01wdktkdWpNdWd3RFFZSktvWklodmNOQVFFTEJRQXdiREVMTUFrR0ExVUUKQmhNQ1EwNHhFREFPQmdOVkJBZ1RCMEpsYVdwcGJtY3hFREFPQmdOVkJBY1RCMEpsYVdwcGJtY3hFVEFQQmdOVgpCQW9UQ0d4cGRHVnJkV0psTVF3d0NnWURWUVFMRXdOMmNHNHhHREFXQmdOVkJBTU1EMnhwZEdWcmRXSmxYM1p3CmJsOWpZVEFlRncweU1qQTBNVFl3TmpVMU1EQmFGdzB5TXpBME1qQXdORFF4TURkYU1ESXhHakFZQmdOVkJBb1QKRVd4cGRHVnJkV0psTFhad2JpMW5jbkJqTVJRd0VnWURWUVFERXd0bmNuQmpMV05zYVdWdWREQlpNQk1HQnlxRwpTTTQ5QWdFR0NDcUdTTTQ5QXdFSEEwSUFCTjB5eDgydEd6QTdhT1lXdUV6YW8zcC83cWxENURSZWtoU2huWUc4CnNwbm10d0U0WU16TkFGZFRYclNleEhuNEZqYjd5cWxtdmNsT2d0ZEtVaC95eXF5alNEQkdNQTRHQTFVZER3RUIKL3dRRUF3SUZvREFUQmdOVkhTVUVEREFLQmdnckJnRUZCUWNEQWpBZkJnTlZIU01FR0RBV2dCVHlkbDJsdmdrYgpTZzJndnZ2N2xaUTNQZVVKQlRBTkJna3Foa2lHOXcwQkFRc0ZBQU9DQVFFQXIyakpXcFZvTGJUSGxZcyszWlpXCmtGVENGL1JNV1RTSCtkaTZVMjNUWWRmY28xSnJQeWFEUkliaWI1QVBLNDZDay85dmdjK1AwMjlXUFBPYzJ5ZjcKMm5uR3IxcDQxVml1VnJjaExvTGd3Sk1uOFhibjBpTUQ1cGdnMERwcVVmVGhqN0ZwWjZuUC9SUkpnVElWS1psQwo0WEhrMW5vVm40b0M5UW0wNGkrZmFoMGVaeFl5ay9udUxsNjVjeEdUd1JZZHAxQ21VdUIxSmJCa2FEYnM4UmpICjk1VUJXTDBFK2taUm45UFhiZHhjYnNsUStTbzFOZElybFU2RndjcVdHak5DekVtTFNQOFhpZnZETmJtbWpkUXMKVnZoTjVZUmhHRHpwd0EyNEJlN1FFUHpIaHlncjJENFowWGI3Z2ZyekRzK0lzdnpvYlEvb2YwM1VHdHV1V0pNUApiQT09Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K",
  "NetworkCaCert": "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURxRENDQXBDZ0F3SUJBZ0lVYW1GNkF6L0ZCSmlrOHVRQUFuV1dvdWRFOWVJd0RRWUpLb1pJaHZjTkFRRUwKQlFBd2JERUxNQWtHQTFVRUJoTUNRMDR4RURBT0JnTlZCQWdUQjBKbGFXcHBibWN4RURBT0JnTlZCQWNUQjBKbAphV3BwYm1jeEVUQVBCZ05WQkFvVENHeHBkR1ZyZFdKbE1Rd3dDZ1lEVlFRTEV3TjJjRzR4R0RBV0JnTlZCQU1NCkQyeHBkR1ZyZFdKbFgzWndibDlqWVRBZUZ3MHlNakEwTVRZd05qVTFNREJhRncweU56QTBNVFV3TmpVMU1EQmEKTUd3eEN6QUpCZ05WQkFZVEFrTk9NUkF3RGdZRFZRUUlFd2RDWldscWFXNW5NUkF3RGdZRFZRUUhFd2RDWldscQphVzVuTVJFd0R3WURWUVFLRXdoc2FYUmxhM1ZpWlRFTU1Bb0dBMVVFQ3hNRGRuQnVNUmd3RmdZRFZRUUREQTlzCmFYUmxhM1ZpWlY5MmNHNWZZMkV3Z2dFaU1BMEdDU3FHU0liM0RRRUJBUVVBQTRJQkR3QXdnZ0VLQW9JQkFRRE8KakRHQ20vRUdLTFJoZ0tsMU8xeVEwamEvYlhjN1VpQnBoNzh6Sm92aStKanZYbVhTdGJpTElxd1o1YkdiMEdKWApLUW13bmRDVmgycUpFUXZCZ2FGc0pMdTFRY09uMU5TZ3plV3ZtYk9yNHhJQ0x1Y2QyRVZWWForNXplNmxSSjcxCjh3aHRpYjVoSlZqdTVWY1VVY2dJNStOWm5zSWdRQ2htaDcvM2RXRkh0MmR3TGxyaVpFcVVXTldLemRPeDI0MFYKeXI4QmY5UXg4MFRhYmFDVHZTUmx0TkE0TFdTYXRMTFZlV2JnMkRDMmQvcUwyMEhaQ1pPek9kT1FDc29xL3d5TwpURkZQL2FHUWFUU0lYRWZLWWtwR01oV1V2RlNoY1pObzBkcEg4czU4UkZweU5wUzJCNXcwV0NyRHRHTlVMN2tQCm1IWlBOVC9GVHhLVDNPcStJdTkvQWdNQkFBR2pRakJBTUE0R0ExVWREd0VCL3dRRUF3SUJCakFQQmdOVkhSTUIKQWY4RUJUQURBUUgvTUIwR0ExVWREZ1FXQkJUeWRsMmx2Z2tiU2cyZ3Z2djdsWlEzUGVVSkJUQU5CZ2txaGtpRwo5dzBCQVFzRkFBT0NBUUVBcnhIZk15YnpLV05jcDQ2MWppOGVWbTNqWXVldUg5ZmRXTUc2ZUxPYTFtNUZJdk5ZCm5LWUVCS1pUTEV4T1N0SklHbUZqbGpGZjNPaEhaSllmcDZnUGFFNGVuUEdDQ1Q1Wkhlcnd3dnNOL2VIVXcwR2wKWlZTSTR1OGozN3ZLZVd4aWZTaHFvc09lWHNhdW5sanlwSVQrWDBQSEx0SnRpMS9pbHlRcmVZaWNTNkIwMGhoKwpCNW1HNUpKRCs5WnVnL2hPbkQ5WXI4RU5mdnNKQzcxemdVVUxGV1NiVi96T1ZhNHlJdzhRU0pnaHFRSmR5MkpRClhGSXp6aXhseGU3NWR1UUVaNGlIWUZhblVHNVlPMFFsZ2FPazZocUp2aU1uUEhaTWorSEFPT2RPTEEwQk15RkcKWnJDMit6b2FMSERrTlExQTdQRHFoMk1VMUFOcGRqU0EvQWhYN0E9PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==",
  "NetworkClientKey": "LS0tLS1CRUdJTiBFQyBQUklWQVRFIEtFWS0tLS0tCk1IY0NBUUVFSUFKbERlbDI3V0h5VjVhbXZ2RzRTSDBIMGx4QXhpVDE3a0lodDJWTjRrUHFvQW9HQ0NxR1NNNDkKQXdFSG9VUURRZ0FFT0JaOUdqVkc2RVpOYzNEMDA4bHJUa0FoaWNHbC81RXZMOTVqWWtPMyt0QnY3UEI4cm1IQQo2eEdzWE96d3lmZEVjOU8vSW5FMWVMRFhZTmE4SURHNFJnPT0KLS0tLS1FTkQgRUMgUFJJVkFURSBLRVktLS0tLQo=",
  "NetworkClientCert": "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNsekNDQVgrZ0F3SUJBZ0lJZlVjckNrM1JGOWt3RFFZSktvWklodmNOQVFFTEJRQXdiREVMTUFrR0ExVUUKQmhNQ1EwNHhFREFPQmdOVkJBZ1RCMEpsYVdwcGJtY3hFREFPQmdOVkJBY1RCMEpsYVdwcGJtY3hFVEFQQmdOVgpCQW9UQ0d4cGRHVnJkV0psTVF3d0NnWURWUVFMRXdOMmNHNHhHREFXQmdOVkJBTU1EMnhwZEdWcmRXSmxYM1p3CmJsOWpZVEFlRncweU1qQTBNVFl3TmpVMU1EQmFGdzB5TXpBME1qQXdORFF4TURkYU1Dd3hGVEFUQmdOVkJBb1QKREd4cGRHVnJkV0psTFhad2JqRVRNQkVHQTFVRUF4TUtkbkJ1TFdOc2FXVnVkREJaTUJNR0J5cUdTTTQ5QWdFRwpDQ3FHU000OUF3RUhBMElBQkRnV2ZSbzFSdWhHVFhOdzlOUEphMDVBSVluQnBmK1JMeS9lWTJKRHQvclFiK3p3CmZLNWh3T3NSckZ6czhNbjNSSFBUdnlKeE5YaXcxMkRXdkNBeHVFYWpTREJHTUE0R0ExVWREd0VCL3dRRUF3SUYKb0RBVEJnTlZIU1VFRERBS0JnZ3JCZ0VGQlFjREFqQWZCZ05WSFNNRUdEQVdnQlR5ZGwybHZna2JTZzJndnZ2NwpsWlEzUGVVSkJUQU5CZ2txaGtpRzl3MEJBUXNGQUFPQ0FRRUFEV3h3VjRzeDR5ZGNVd2NMSWZJVUV3UEF1SkxLCnJMd29BNlAvSFZ3WTk1VWNUTlVlTXV1RTJyZExsMXBqaGtjU09tT0pNRENsQ3hyMlpTRzAvZ09Ha0pHNEVIQjIKMkhnNE91djJtSEloWmVtampuNWZIY0Z3TkJrQjYyU21Tem8vQWV4eUp1N1lyeUErUnBqaUFmWHl6WjVDaU5NdwoyVyt1RXBsN05adG1TRUxCRC85b0F6dDFzTjZUUWV3K3k0ZXlQdGh4Z09Wc08xcXVRZm5LcGNzVVJjaUtvQ0hzCjhIbnZ5SVNrdDNiTnEvZnp0UGd2VXBMMkJZNDJJR1pVcFFZVG5JNFJDeStVMThMTGU3TjdHa0VhYXp1VVBVMnoKUzZBMXJsZ1N0NDl5THlpRTRMUUtBcUdhYnZ5bEs5MFMvK01lRXpnNSsyTG92NkUydE5WOGxSWW9Gdz09Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K"
}
```

#### 节点取消注册 UnRegister

- demo：取消node-token为b52f93d3f0ec4be7的注册绑定

```shell
grpcurl -d '{"token": "b52f93d3f0ec4be7"}' -cacert ca.pem -cert client.pem -key client-key.pem 101.43.253.110:6440 pb.LiteKubeNCService.CheckConnState
```

- UnRegisterRequest参数：空参

| **参数** | 类型   | 含义       | 是否必须 | demo             |
| -------- | ------ | ---------- | -------- | ---------------- |
| token    | string | node-token | 是       | b52f93d3f0ec4be7 |

- UnRegisterResponse返回数据

```json
{
  "code": "200",
  "message": "ok",
  "result": true
}
```

#### 检查连接状态 CheckConnState

- demo：获取node-token为b52f93d3f0ec4be7的连接状态

```shell
grpcurl -d '{"token": "b52f93d3f0ec4be7"}' -cacert ca.pem -cert client.pem -key client-key.pem 101.43.253.110:6440 pb.LiteKubeNCService.CheckConnState
```

- CheckConnStateRequest参数

| **参数** | 类型   | 含义       | 是否必须 | demo             |
| -------- | ------ | ---------- | -------- | ---------------- |
| token    | string | node-token | 是       | b52f93d3f0ec4be7 |

- CheckConnStateResponse返回数据

```json
{
  "code": "200",
  "message": "ok",
  "connState": 3,
  "bindIp": "10.1.1.3"
}
```

> connState状态码

| connState枚举类型 | 值   | 含义           |
| ----------------- | ---- | -------------- |
| STATE_IDLE        | -1   | 断连状态       |
| STATE_INIT        | 1    | 初始化连接阶段 |
| STATE_CONNECTED   | 3    | 保持连接       |

#### 获取连接节点ip GetRegistedIp

可以复用CheckConnState接口

