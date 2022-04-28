package grpc_client

import (
	"context"
	"encoding/base64"
	"github.com/Litekube/network-controller/grpc/pb_gen"
	certutil "github.com/rancher/dynamiclistener/cert"
	"testing"
	"time"
)

func TestGetGrpcClient(t *testing.T) {
	// from grpc server: ca2.pem client.pem client-key.pem
	client := &GrpcClient{
		Ip:          "101.43.253.110",
		Port:        "6440",
		GrpcCertDir: "/Users/zhujianxing/GoLandProjects/network-controller/certs/test2/",
		CAFile:      "ca2.pem",
		CertFile:    "client.pem",
		KeyFile:     "client-key.pem",
	}
	err := client.InitGrpcClientConn()
	logger.Info(err)
	req := &pb_gen.HelloWorldRequest{HelloText: "hello~"}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resp, err := client.C.HelloWorld(ctx, req)
	logger.Info(resp)
	logger.Info(err)
}

func TestGrpcClient_InitGrpcClientConn(t *testing.T) {
	client := &GrpcBootStrapClient{
		Ip:            "101.43.253.110",
		BootstrapPort: "6439",
	}
	err := client.InitGrpcBootstrapClientConn()
	logger.Info(err)
	req := &pb_gen.GetTokenRequest{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resp, err := client.BootstrapC.GetToken(ctx, req)
	logger.Info(resp)
	logger.Info(err)

	caBytes, err := base64.StdEncoding.DecodeString(resp.GrpcCaCert)
	certBytes, err := base64.StdEncoding.DecodeString(resp.GrpcClientCert)
	keyBytes, err := base64.StdEncoding.DecodeString(resp.GrpcClientKey)
	certutil.WriteCert("/Users/zhujianxing/GoLandProjects/network-controller/certs/test/ca2.pem", caBytes)
	certutil.WriteCert("/Users/zhujianxing/GoLandProjects/network-controller/certs/test/client.pem", certBytes)
	certutil.WriteKey("/Users/zhujianxing/GoLandProjects/network-controller/certs/test/client-key.pem", keyBytes)
}

// test server self-gen cert(ca/server/client)
func TestGrpcClient_InitGrpcClientConn2(t *testing.T) {
	client := &GrpcClient{
		Ip:          "101.43.253.110",
		Port:        "6440",
		GrpcCertDir: "/Users/zhujianxing/GoLandProjects/network-controller/certs/test1/",
		CAFile:      "ca.pem",
		CertFile:    "client.pem",
		KeyFile:     "client-key.pem",
	}
	err := client.InitGrpcClientConn()
	logger.Info(err)
	req := &pb_gen.HelloWorldRequest{}
	//md := metadata.New(map[string]string{
	//	"node-token": "xx",
	//})
	//ctx := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := client.C.HelloWorld(context.Background(), req)
	logger.Info(resp)
	logger.Info(err)
}
