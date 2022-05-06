package main

import (
	"encoding/base64"
	"fmt"
	certutil "github.com/rancher/dynamiclistener/cert"
	"testing"
)

func TestGetBootstrapToken(t *testing.T) {
	Init()
	resp, err := GetBootstrapToken()
	if err != nil {
		fmt.Printf("fail to call GetBootstrapToken err: %+v\n", err)
		return
	}
	fmt.Printf("get bootstrap token:%+v\n", resp.BootStrapToken)
}

func TestGetToken(t *testing.T) {
	InitBootstrapClient()
	tokenResp, err := GetToken("1b678226bb104be8")
	if err != nil {
		fmt.Printf("fail to call GetToken err: %+v\n", err)
		return
	}
	token := tokenResp.Token
	fmt.Printf("register token: %+v\n", token)

	caBytes, err := base64.StdEncoding.DecodeString(tokenResp.GrpcCaCert)
	certBytes, err := base64.StdEncoding.DecodeString(tokenResp.GrpcClientCert)
	keyBytes, err := base64.StdEncoding.DecodeString(tokenResp.GrpcClientKey)
	certutil.WriteCert("/Users/zhujianxing/GoLandProjects/network-controller/certs/test1/ca.pem", caBytes)
	certutil.WriteCert("/Users/zhujianxing/GoLandProjects/network-controller/certs/test1/client.pem", certBytes)
	certutil.WriteKey("/Users/zhujianxing/GoLandProjects/network-controller/certs/test1/client-key.pem", keyBytes)

	caBytes, err = base64.StdEncoding.DecodeString(tokenResp.NetworkCaCert)
	certBytes, err = base64.StdEncoding.DecodeString(tokenResp.NetworkClientCert)
	keyBytes, err = base64.StdEncoding.DecodeString(tokenResp.NetworkClientKey)
	certutil.WriteCert("/Users/zhujianxing/GoLandProjects/network-controller/certs/test2/ca.pem", caBytes)
	certutil.WriteCert("/Users/zhujianxing/GoLandProjects/network-controller/certs/test2/client.pem", certBytes)
	certutil.WriteKey("/Users/zhujianxing/GoLandProjects/network-controller/certs/test2/client-key.pem", keyBytes)

}

func TestCheckConnState(t *testing.T) {
	Init()
	checkResp, err := CheckConnState("009794b89caa4881")
	if err != nil {
		fmt.Printf("fail to call CheckConnState err: %+v\n", err)
		return
	}
	fmt.Printf("get bind ip:%+v, conn state:%+v\n", checkResp.BindIp, checkResp.ConnState)
}

func TestUnRegister(t *testing.T) {
	Init()
	unRegisResp, err := UnRegister("e95b4398c1514a24")
	if err != nil {
		fmt.Printf("fail to call UnRegister err: %+v\n", err)
		return
	}
	fmt.Printf("if succeed to unRegister: %+v\n", unRegisResp.Result)
}
