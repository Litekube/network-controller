package main

import (
	"context"
	"github.com/Litekube/network-controller/grpc/grpc_client"
	"github.com/Litekube/network-controller/grpc/pb_gen"
	"github.com/Litekube/network-controller/utils"
)

var logger = utils.GetLogger()

var Client *grpc_client.GrpcClient

func Init() {
	Client = &grpc_client.GrpcClient{
		Ip:          "101.43.253.110",
		Port:        "6440",
		GrpcCertDir: "/root/go_project/network-controller-usage/certs/test1/",
		CAFile:      "ca.pem",
		CertFile:    "client.pem",
		KeyFile:     "client-key.pem",
	}
	err := Client.InitGrpcClientConn()
	logger.Info(err)
}

func Init2() {
	Client = &grpc_client.GrpcClient{
		Ip:   "101.43.253.110",
		Port: "6440",
	}
	err := Client.InitGrpcClientConn()
	logger.Info(err)
}

func main() {
	Init()
}

func GetBootstrapToken() (*pb_gen.GetBootStrapTokenResponse, error) {
	req := &pb_gen.GetBootStrapTokenRequest{
		ExpireTime: 3,
	}

	resp, err := Client.C.GetBootStrapToken(context.Background(), req)
	logger.Info(resp)
	logger.Info(err)

	return resp, nil
}

func GetToken(bootstrapToken string) (*pb_gen.GetTokenResponse, error) {
	req := &pb_gen.GetTokenRequest{
		BootStrapToken: bootstrapToken,
	}

	resp, err := Client.C.GetToken(context.Background(), req)
	logger.Info(resp)
	logger.Info(err)

	return resp, nil
}

func CheckConnState(token string) (*pb_gen.CheckConnResponse, error) {
	req := &pb_gen.CheckConnStateRequest{
		Token: token,
	}

	resp, err := Client.C.CheckConnState(context.Background(), req)
	logger.Info(resp)
	logger.Info(err)
	return resp, err
}

func UnRegister(token string) (*pb_gen.UnRegisterResponse, error) {
	req := &pb_gen.UnRegisterRequest{
		Token: token,
	}

	resp, err := Client.C.UnRegister(context.Background(), req)
	logger.Info(resp)
	logger.Info(err)
	return resp, err
}
