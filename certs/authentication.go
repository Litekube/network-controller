package certs

import (
	"github.com/litekube/LiteKube/pkg/global"
	"github.com/rancher/dynamiclistener/cert"
	"litekube-vpn/config"
	"litekube-vpn/utils"
	"net"
)

func CheckGrpcCertConfig(tlsConfig config.TLSConfig) error {
	// generate for grpc
	// generate CA
	regenGrpc, err := GenerateSigningCertKey(false, "litekube-vpn-grpc", tlsConfig.CAFile, tlsConfig.CAKeyFile)
	if err != nil {
		return err
	}

	// generate server
	if _, _, _, err := GenerateServerCertKey(regenGrpc, "litekube-vpn-grpc-server", nil,
		&cert.AltNames{
			DNSNames: append([]string{}, global.LocalHostDNSName),
			IPs:      append(append(global.LocalIPs, []net.IP{net.ParseIP(utils.QueryPublicIp())}...)),
		}, tlsConfig.CAFile, tlsConfig.CAKeyFile, tlsConfig.ServerCertFile, tlsConfig.ServerKeyFile); err != nil {
		return err
	}

	// generate client
	if _, _, _, err := GenerateClientCertKey(regenGrpc, "litekube-vpn-grpc-client", []string{"litekube-vpn-grpc"}, tlsConfig.CAFile, tlsConfig.CAKeyFile, tlsConfig.ClientCertFile, tlsConfig.ClientKeyFile); err != nil {
		return err
	}
	return nil
}

func CheckVpnCertConfig(tlsConfig config.TLSConfig) error {
	//generate for vpn
	//generate CA
	regenGrpc, err := GenerateSigningCertKey(false, "litekube-vpn", tlsConfig.CAFile, tlsConfig.CAKeyFile)
	if err != nil {
		return err
	}

	// generate server
	if _, _, _, err := GenerateServerCertKey(regenGrpc, "litekube-vpn-server", nil,
		&cert.AltNames{
			DNSNames: append([]string{}, global.LocalHostDNSName),
			IPs:      append(append(global.LocalIPs, []net.IP{net.ParseIP(utils.QueryPublicIp())}...)),
		}, tlsConfig.CAFile, tlsConfig.CAKeyFile, tlsConfig.ServerCertFile, tlsConfig.ServerKeyFile); err != nil {
		return err
	}

	// generate client
	if _, _, _, err := GenerateClientCertKey(regenGrpc, "litekube-vpn-client", []string{"litekube-vpn"}, tlsConfig.CAFile, tlsConfig.CAKeyFile, tlsConfig.ClientCertFile, tlsConfig.ClientKeyFile); err != nil {
		return err
	}
	return nil
}
