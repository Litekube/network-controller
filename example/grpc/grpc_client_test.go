package main

import (
	"encoding/base64"
	certutil "github.com/rancher/dynamiclistener/cert"
	"testing"
)

func TestGetBootstrapToken(t *testing.T) {
	Init()
	resp, err := GetBootstrapToken()
	if err != nil {
		logger.Errorf("fail to call GetBootstrapToken err: %+v", err)
		return
	}
	logger.Infof("get bootstrap token:%+v", resp.BootStrapToken)
}

func TestGetToken(t *testing.T) {
	Init2()
	tokenResp, err := GetToken("70318e4cb8dd4100")
	if err != nil {
		logger.Errorf("fail to call GetToken err: %+v", err)
		return
	}
	token := tokenResp.Token
	logger.Infof("register token: %+v", token)

	caBytes, err := base64.StdEncoding.DecodeString(tokenResp.GrpcCaCert)
	certBytes, err := base64.StdEncoding.DecodeString(tokenResp.GrpcClientCert)
	keyBytes, err := base64.StdEncoding.DecodeString(tokenResp.GrpcClientKey)
	certutil.WriteCert("/root/go_project/network-controller-usage/certs/test1/ca.pem", caBytes)
	certutil.WriteCert("/root/go_project/network-controller-usage/certs/test1/client.pem", certBytes)
	certutil.WriteKey("/root/go_project/network-controller-usage/certs/test1/client-key.pem", keyBytes)

	caBytes, err = base64.StdEncoding.DecodeString(tokenResp.NetworkCaCert)
	certBytes, err = base64.StdEncoding.DecodeString(tokenResp.NetworkClientCert)
	keyBytes, err = base64.StdEncoding.DecodeString(tokenResp.NetworkClientKey)
	certutil.WriteCert("/root/go_project/network-controller-usage/certs/test2/ca.pem", caBytes)
	certutil.WriteCert("/root/go_project/network-controller-usage/certs/test2/client.pem", certBytes)
	certutil.WriteKey("/root/go_project/network-controller-usage/certs/test2/client-key.pem", keyBytes)

}

func TestCheckConnState(t *testing.T) {
	Init()
	checkResp, err := CheckConnState("e95b4398c1514a24")
	if err != nil {
		logger.Errorf("fail to call CheckConnState err: %+v", err)
		return
	}
	logger.Infof("get bind ip:%+v, conn state:%+v", checkResp.BindIp, checkResp.ConnState)
}

func TestUnRegister(t *testing.T) {
	Init()
	unRegisResp, err := UnRegister("e95b4398c1514a24")
	if err != nil {
		logger.Errorf("fail to call UnRegister err: %+v", err)
		return
	}
	logger.Infof("if succeed to unRegister: %+v", unRegisResp.Result)
}
