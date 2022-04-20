package internal

import (
	"context"
	"encoding/base64"
	"errors"
	certutil "github.com/rancher/dynamiclistener/cert"
	"litekube-vpn/certs"
	"litekube-vpn/config"
	"litekube-vpn/contant"
	"litekube-vpn/grpc/pb_gen"
	"litekube-vpn/sqlite"
	"litekube-vpn/utils"
)

type LiteVpnService struct {
	unRegisterCh  chan string
	grpcTlsConfig config.TLSConfig
	vpnTlsConfig  config.TLSConfig
}

var logger = utils.GetLogger()

func NewLiteVpnService(unRegisterCh chan string, grpcTlsConfig config.TLSConfig, vpnTlsConfig config.TLSConfig) *LiteVpnService {
	return &LiteVpnService{
		unRegisterCh:  unRegisterCh,
		grpcTlsConfig: grpcTlsConfig,
		vpnTlsConfig:  vpnTlsConfig,
	}
}

func (service *LiteVpnService) CheckConnState(ctx context.Context, req *pb_gen.CheckConnStateRequest) (*pb_gen.CheckConnResponse, error) {

	wrappedResp := func(code, message, bindIp string, state int32) (resp *pb_gen.CheckConnResponse, err error) {
		if code != contant.STATUS_OK {
			logger.Errorf("query token: %+v err: %+v", req.Token, err)
			err = errors.New(message)
		}
		resp = &pb_gen.CheckConnResponse{
			Message:   message,
			Code:      code,
			ConnState: state,
			BindIp:    bindIp,
		}
		logger.Debugf("resp: %+v", resp)
		return
	}

	if len(req.Token) == 0 {
		return wrappedResp(contant.STATUS_BADREQUEST, "token can't be empty", "", -1)
	}

	vpnMgr := sqlite.VpnMgr{}
	item, err := vpnMgr.QueryByToken(req.Token)
	if item == nil {
		return wrappedResp(contant.STATUS_OK, err.Error(), "", -1)
	} else if err != nil {
		return wrappedResp(contant.STATUS_ERR, err.Error(), "", -1)
	}

	return wrappedResp(contant.STATUS_OK, contant.MESSAGE_OK, item.BindIp, int32(item.State))
}

func (service *LiteVpnService) UnRegister(ctx context.Context, req *pb_gen.UnRegisterRequest) (*pb_gen.UnRegisterResponse, error) {

	wrappedResp := func(code, message string, result bool) (resp *pb_gen.UnRegisterResponse, err error) {
		if code != contant.STATUS_OK {
			logger.Errorf("query token: %+v err: %+v", req.Token, err)
			err = errors.New(message)
		}
		resp = &pb_gen.UnRegisterResponse{
			Message: message,
			Code:    code,
			Result:  result,
		}
		logger.Debugf("resp: %+v", resp)
		return
	}

	if len(req.Token) == 0 {
		return wrappedResp(contant.STATUS_BADREQUEST, "token can't be empty", false)
	}

	vpnMgr := sqlite.VpnMgr{}
	item, err := vpnMgr.QueryByToken(req.Token)
	if item == nil {
		return wrappedResp(contant.STATUS_OK, err.Error(), false)
	} else if err != nil {
		return wrappedResp(contant.STATUS_ERR, err.Error(), false)
	}

	result, err := vpnMgr.DeleteById(item.Id)
	if err != nil {
		return wrappedResp(contant.STATUS_ERR, err.Error(), result)
	}

	service.unRegisterCh <- item.BindIp
	return wrappedResp(contant.STATUS_OK, contant.MESSAGE_OK, result)
}

func (service *LiteVpnService) GetRegistedIp(ctx context.Context, req *pb_gen.GetRegistedIpRequest) (*pb_gen.GetRegistedIpResponse, error) {

	wrappedResp := func(code, message, ip string) (resp *pb_gen.GetRegistedIpResponse, err error) {
		if code != contant.STATUS_OK {
			logger.Errorf("query token: %+v err: %+v", req.Token, err)
			err = errors.New(message)
		}
		resp = &pb_gen.GetRegistedIpResponse{
			Message: message,
			Code:    code,
			Ip:      ip,
		}
		logger.Debugf("resp: %+v", resp)
		return
	}

	if len(req.Token) == 0 {
		return wrappedResp(contant.STATUS_BADREQUEST, "token can't be empty", "")
	}

	vpnMgr := sqlite.VpnMgr{}
	item, err := vpnMgr.QueryByToken(req.Token)
	if item == nil {
		return wrappedResp(contant.STATUS_OK, err.Error(), "")
	} else if err != nil {
		return wrappedResp(contant.STATUS_ERR, err.Error(), "")
	}

	return wrappedResp(contant.STATUS_OK, contant.MESSAGE_OK, item.BindIp)
}

func (service *LiteVpnService) GetToken(ctx context.Context, req *pb_gen.GetTokenRequest) (*pb_gen.GetTokenResponse, error) {

	wrappedResp := func(code, message, token string) (resp *pb_gen.GetTokenResponse, err error) {
		if code != contant.STATUS_OK {
			err = errors.New(message)
		}
		resp = &pb_gen.GetTokenResponse{
			Code:           code,
			Message:        message,
			Token:          token,
			GrpcCaCert:     "",
			GrpcClientKey:  "",
			GrpcClientCert: "",
			VpnCaCert:      "",
			VpnClientKey:   "",
			VpnClientCert:  "",
		}
		logger.Debugf("resp: %+v", resp)
		return
	}

	token := utils.GetUniqueToken()
	vpnMgr := sqlite.VpnMgr{}
	// no need
	//item, err := vpnMgr.QueryByToken(token)
	err := vpnMgr.Insert(sqlite.VpnMgr{
		Token:  token,
		State:  contant.STATE_IDLE,
		BindIp: "",
	})
	if err != nil {
		return wrappedResp(contant.STATUS_ERR, err.Error(), "")
	}

	keyBytes, certBytes, _, err := certs.GenerateClientCertKey(true, "litekube-vpn-grpc-client", []string{"litekube-vpn-grpc"}, service.grpcTlsConfig.CAFile, service.grpcTlsConfig.CAKeyFile, service.grpcTlsConfig.ClientCertFile, service.grpcTlsConfig.ClientKeyFile)
	if err != nil {
		return wrappedResp(contant.STATUS_ERR, err.Error(), "")
	}

	resp, _ := wrappedResp(contant.STATUS_OK, contant.MESSAGE_OK, token)

	// load grpc ca.pem client.pem client-key.pem
	grpcCaCert, err := certs.LoadCertificate(service.grpcTlsConfig.CAFile)
	if err != nil {
		return wrappedResp(contant.STATUS_ERR, err.Error(), "")
	}
	resp.GrpcCaCert = base64.StdEncoding.EncodeToString(certutil.EncodeCertPEM(grpcCaCert))
	resp.GrpcClientKey = base64.StdEncoding.EncodeToString(keyBytes)
	resp.GrpcClientCert = base64.StdEncoding.EncodeToString(certBytes)

	keyBytes, certBytes, _, err = certs.GenerateClientCertKey(true, "litekube-vpn-client", []string{"litekube-vpn"}, service.vpnTlsConfig.CAFile, service.vpnTlsConfig.CAKeyFile, service.vpnTlsConfig.ClientCertFile, service.vpnTlsConfig.ClientKeyFile)
	if err != nil {
		return wrappedResp(contant.STATUS_ERR, err.Error(), "")
	}

	// load vpn ca.pem client.pem client-key.pem
	vpnCaCert, err := certs.LoadCertificate(service.vpnTlsConfig.CAFile)
	if err != nil {
		return wrappedResp(contant.STATUS_ERR, err.Error(), "")
	}
	resp.VpnCaCert = base64.StdEncoding.EncodeToString(certutil.EncodeCertPEM(vpnCaCert))
	resp.VpnClientKey = base64.StdEncoding.EncodeToString(keyBytes)
	resp.VpnClientCert = base64.StdEncoding.EncodeToString(certBytes)
	return resp, nil
}
