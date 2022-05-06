package main

import (
	"context"
	"fmt"
	"github.com/Litekube/network-controller/grpc/grpc_client"
	"github.com/Litekube/network-controller/grpc/pb_gen"
)

var Client *grpc_client.GrpcClient
var BootstrapClient *grpc_client.GrpcBootStrapClient

func Init() {
	Client = &grpc_client.GrpcClient{
		Ip:          "127.0.0.1",
		Port:        "6440",
		GrpcCertDir: "/root/go_project/network-controller/certs/grpc/",
		CAFile:      "ca.pem",
		CertFile:    "client.pem",
		KeyFile:     "client-key.pem",
	}
	err := Client.InitGrpcClientConn()
	fmt.Println(err)
}

func InitBootstrapClient() {
	BootstrapClient = &grpc_client.GrpcBootStrapClient{
		Ip:            "101.43.253.110",
		BootstrapPort: "6439",
	}
	err := BootstrapClient.InitGrpcBootstrapClientConn()
	fmt.Println(err)
}

func main() {
	Init()
	InitBootstrapClient()
}

func GetBootstrapToken() (*pb_gen.GetBootStrapTokenResponse, error) {
	req := &pb_gen.GetBootStrapTokenRequest{
		ExpireTime: 3,
	}

	resp, err := Client.C.GetBootStrapToken(context.Background(), req)
	fmt.Println(resp)
	fmt.Println(err)

	return resp, err
}

func GetToken(bootstrapToken string) (*pb_gen.GetTokenResponse, error) {
	req := &pb_gen.GetTokenRequest{
		BootStrapToken: bootstrapToken,
	}

	resp, err := BootstrapClient.BootstrapC.GetToken(context.Background(), req)
	fmt.Println(resp)
	fmt.Println(err)

	return resp, err
}

func CheckConnState(token string) (*pb_gen.CheckConnResponse, error) {
	req := &pb_gen.CheckConnStateRequest{
		Token: token,
	}

	resp, err := Client.C.CheckConnState(context.Background(), req)
	fmt.Println(resp)
	fmt.Println(err)
	return resp, err
}

func UnRegister(token string) (*pb_gen.UnRegisterResponse, error) {
	req := &pb_gen.UnRegisterRequest{
		Token: token,
	}

	resp, err := Client.C.UnRegister(context.Background(), req)
	fmt.Println(resp)
	fmt.Println(err)
	return resp, err
}
