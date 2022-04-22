package grpc_client

import (
	"context"
	"encoding/base64"
	"github.com/Litekube/litekube-vpn/grpc/pb_gen"
	certutil "github.com/rancher/dynamiclistener/cert"
	"testing"
	"time"
)

func TestGetGrpcClient(t *testing.T) {
	// from grpc server: ca2.pem client.pem client-key.pem
	client := &GrpcClient{
		Ip:          "101.43.253.110",
		Port:        "6440",
		grpcCertDir: "/Users/zhujianxing/GoLandProjects/litekube-vpn/certs/test/",
		CAFile:      "ca2.pem",
		CertFile:    "client.pem",
		KeyFile:     "client-key.pem",
	}
	err := client.InitGrpcClientConn()
	logger.Info(err)
	req := &pb_gen.HelloWorldRequest{HelloText: "hello~"}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resp, err := client.c.HelloWorld(ctx, req)
	logger.Info(resp)
	logger.Info(err)
}

func TestGrpcClient_InitGrpcClientConn(t *testing.T) {
	client := &GrpcClient{
		Ip:          "101.43.253.110",
		Port:        "6440",
		grpcCertDir: "/Users/zhujianxing/GoLandProjects/litekube-vpn/certs/test1/",
		CAFile:      "ca.pem",
		CertFile:    "client.pem",
		KeyFile:     "client-key.pem",
	}
	err := client.InitGrpcClientConn()
	logger.Info(err)
	req := &pb_gen.GetTokenRequest{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resp, err := client.c.GetToken(ctx, req)
	logger.Info(resp)
	logger.Info(err)

	caBytes, err := base64.StdEncoding.DecodeString(resp.GrpcCaCert)
	certBytes, err := base64.StdEncoding.DecodeString(resp.GrpcClientCert)
	keyBytes, err := base64.StdEncoding.DecodeString(resp.GrpcClientKey)
	certutil.WriteCert("/Users/zhujianxing/GoLandProjects/litekube-vpn/certs/test/ca2.pem", caBytes)
	certutil.WriteCert("/Users/zhujianxing/GoLandProjects/litekube-vpn/certs/test/client.pem", certBytes)
	certutil.WriteKey("/Users/zhujianxing/GoLandProjects/litekube-vpn/certs/test/client-key.pem", keyBytes)
}

// test server self-gen cert(ca/server/client)
func TestGrpcClient_InitGrpcClientConn2(t *testing.T) {
	client := &GrpcClient{
		Ip:          "101.43.253.110",
		Port:        "6440",
		grpcCertDir: "/Users/zhujianxing/GoLandProjects/litekube-vpn/certs/test1/",
		CAFile:      "ca.pem",
		CertFile:    "client.pem",
		KeyFile:     "client-key.pem",
	}
	err := client.InitGrpcClientConn()
	logger.Info(err)
	req := &pb_gen.HelloWorldRequest{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resp, err := client.c.HelloWorld(ctx, req)
	logger.Info(resp)
	logger.Info(err)
}
