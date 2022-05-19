English | [简体中文](./API-explain_CN.md)

## gRPC API Doc

* [gRPC API Doc](#grpc-api-doc)
  * [Overview](#overview)
  * [grpcurl tool](#grpcurl-tool)
  * [APIs](#apis)
    * [Bootstrap Register：GetBootStrapToken](#bootstrap-registergetbootstraptoken)
    * [Node register: GetToken](#node-register-gettoken)
    * [Node unregister: UnRegister](#node-unregister-unregister)
    * [Check node connection state: CheckConnState](#check-node-connection-state-checkconnstate)
    * [Get connected node ip: GetRegistedIp](#get-connected-node-ip-getregistedip)

### Overview

Realize communication interaction service based on tcp gRPC+protobuf, support tls secure communication

```
service LiteKubeNCService {
  rpc GetBootStrapToken(GetBootStrapTokenRequest) returns (GetBootStrapTokenResponse) {}
  rpc GetToken(GetTokenRequest) returns (GetTokenResponse) {}
  rpc CheckConnState(CheckConnStateRequest) returns (CheckConnResponse){}
  rpc UnRegister(UnRegisterRequest) returns (UnRegisterResponse){}
  rpc GetRegistedIp(GetRegistedIpRequest) returns (GetRegistedIpResponse){}
}
```

- status code explaination

| status code       |      | meaning                              |
| ----------------- | ---- | ------------------------------------ |
| STATUS_OK         | 200  | success                              |
| STATUS_BADREQUEST | 400  | the client parameters are not normal |
| STATUS_ERR        | 500  | server Internal logic error          |

### grpcurl tool

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
# Query the Service List
grpcurl -plaintext 101.43.253.110:6440 list pb.LiteKubeNCService

# Query a request
grpcurl -plaintext 101.43.253.110:6440 describe pb.LiteKubeNCService.HelloWorld

# Query Parameters;
grpcurl -plaintext 101.43.253.110:6440 describe pb.HelloWorldRequest 
grpcurl -plaintext 101.43.253.110:6440 describe pb.HelloWorldResponse

# grpc call
grpcurl -d '{"token": "b52f93d3f0ec4be7"}' -plaintext 101.43.253.110:6440 pb.LiteKubeNCService.CheckConnState

grpcurl -d '{"token": "b52f93d3f0ec4be7"}' -plaintext 101.43.253.110:6440 pb.LiteKubeNCService.UnRegister

grpcurl -d '{"token": "b52f93d3f0ec4be7"}' -plaintext 101.43.253.110:6440 pb.LiteKubeNCService.GetRegistedIp

grpcurl -d '{"bootStrapToken": "deac5f329feb4729"}' -plaintext 101.43.253.110:6439 pb.LiteKubeNCBootstrapService.GetToken

grpcurl -d '{"expireTime": 10}' -plaintext 101.43.253.110:6440 pb.LiteKubeNCService.GetBootStrapToken

```

- support tls

```shell
# grpc call
grpcurl -d '' -insecure 101.43.253.110:6439 pb.LiteKubeNCBootstrapService.HealthCheck

grpcurl -d '' -cacert ca.pem -cert client.pem -key client-key.pem 101.43.253.110:6440 pb.LiteKubeNCService.HealthCheck

grpcurl -d '{"token": "b52f93d3f0ec4be7"}' -cacert ca.pem -cert client.pem -key client-key.pem 101.43.253.110:6440 pb.LiteKubeNCService.CheckConnState

grpcurl -d '{"token": "009794b89caa4881"}' -cacert ca.pem -cert client.pem -key client-key.pem 101.43.253.110:6440 pb.LiteKubeNCService.UnRegister

grpcurl -d '{"token": "009794b89caa4881"}' -cacert ca.pem -cert client.pem -key client-key.pem 101.43.253.110:6440 pb.LiteKubeNCService.GetRegistedIp

grpcurl -d '{"bootStrapToken": "deac5f329feb4729"}' -insecure 101.43.253.110:6439 pb.LiteKubeNCBootstrapService.GetToken

grpcurl -d '{"bootStrapToken": "deac5f329feb4729"}' -insecure 127.0.0.1:6439 pb.LiteKubeNCBootstrapService.GetToken

grpcurl -d '{"bootStrapToken": "65bdd99bf8904634"}' -cacert ca.pem -cert client.pem -key client-key.pem 101.43.253.110:6439 pb.LiteKubeNCBootstrapService.GetToken

grpcurl -d '{"expireTime": -1}' -cacert ca.pem -cert client.pem -key client-key.pem  101.43.253.110:6440 pb.LiteKubeNCService.GetBootStrapToken
```

### APIs

#### Bootstrap Register：GetBootStrapToken

- demo：Get a Bootstrap token with an expiration time of 10 minutes

```shell
grpcurl -d '{"expireTime": 10}' -cacert ca.pem -cert client.pem -key client-key.pem  101.43.253.110:6440 pb.LiteKubeNCService.GetBootStrapToken

grpcurl -d '{"expireTime": -1}' -cacert ca.pem -cert client.pem -key client-key.pem  101.43.253.110:6440 pb.LiteKubeNCService.GetBootStrapToken
```

- GetBootStrapTokenRequest pameters

| parameters | type  | meaning     | required | demo                                        |
| ---------- | ----- | ----------- | -------- | ------------------------------------------- |
| expireTime | int32 | expire time | no       | 10（A negative number means no expiration） |

- GetBootStrapTokenResponse data

```json
{
  "code": "200",
  "message": "ok",
  "bootStrapToken": "deac5f329feb4729",
  "cloudIp": "101.43.253.110",
  "port": "6440"
}
```

#### Node register: GetToken

- demo

```shell
grpcurl -d '{"bootStrapToken": "deac5f329feb4729"}' -insecure 101.43.253.110:6439 pb.LiteKubeNCBootstrapService.GetToken

grpcurl -d '{"bootStrapToken": "deac5f329feb4729"}' -cacert ca.pem -cert client.pem -key client-key.pem 101.43.253.110:6439 pb.LiteKubeNCBootstrapService.GetToken
```

- GetTokenRequest parameters

| parameters     | type   | meaning         | required | demo             |
| -------------- | ------ | --------------- | -------- | ---------------- |
| bootStrapToken | string | bootstrap token | yes      | deac5f329feb4729 |

- GetTokenResponse data（The certificate fields are base64 encoded）

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

#### Node unregister: UnRegister

- demo：Cancel the registration binding of node-token as b52f93d3f0ec4be7

```shell
grpcurl -d '{"token": "b52f93d3f0ec4be7"}' -cacert ca.pem -cert client.pem -key client-key.pem 101.43.253.110:6440 pb.LiteKubeNCService.CheckConnState
```

- UnRegisterRequest parameters

| parameters | type   | meaning    | required | demo             |
| ---------- | ------ | ---------- | -------- | ---------------- |
| token      | string | node-token | yes      | b52f93d3f0ec4be7 |

- UnRegisterResponse data

```json
{
  "code": "200",
  "message": "ok",
  "result": true
}
```

#### Check node connection state: CheckConnState

- demo：Get the connection status with node-token as b52f93d3f0ec4be7

```shell
grpcurl -d '{"token": "b52f93d3f0ec4be7"}' -cacert ca.pem -cert client.pem -key client-key.pem 101.43.253.110:6440 pb.LiteKubeNCService.CheckConnState
```

- CheckConnStateRequest parameters

| parameters | type   | meaning    | required | demo             |
| ---------- | ------ | ---------- | -------- | ---------------- |
| token      | string | node-token | yes      | b52f93d3f0ec4be7 |

- CheckConnStateResponse 

```json
{
  "code": "200",
  "message": "ok",
  "connState": 3,
  "bindIp": "10.1.1.3"
}
```

> connState status code

| connState status code | value | meaning                         |
| --------------------- | ----- | ------------------------------- |
| STATE_IDLE            | -1    | close                           |
| STATE_INIT            | 1     | init connect, not connected yet |
| STATE_CONNECTED       | 3     | connected                       |

#### Get connected node ip: GetRegistedIp

CheckConnState interface can be reused

