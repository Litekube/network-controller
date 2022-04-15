package internal

import (
	"context"
	"errors"
	"ws-vpn/grpc/pb_gen"
	"ws-vpn/sqlite"
	"ws-vpn/utils"
	"ws-vpn/vpn"
)

type LiteVpnService struct {
	unRegisterCh chan string
}

var logger = utils.GetLogger()

func NewLiteVpnService(unRegisterCh chan string) *LiteVpnService {
	return &LiteVpnService{
		unRegisterCh: unRegisterCh,
	}
}

func (service *LiteVpnService) CheckConnState(ctx context.Context, req *pb_gen.CheckConnStateRequest) (*pb_gen.CheckConnResponse, error) {

	wrappedResp := func(code, message, bindIp string, state int32) (resp *pb_gen.CheckConnResponse, err error) {
		if code != vpn.STATUS_OK {
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
		return wrappedResp(vpn.STATUS_BADREQUEST, "token can't be empty", "", -1)
	}

	vpnMgr := sqlite.VpnMgr{}
	item, err := vpnMgr.QueryByToken(req.Token)
	if item == nil {
		return wrappedResp(vpn.STATUS_OK, err.Error(), "", -1)
	} else if err != nil {
		return wrappedResp(vpn.STATUS_ERR, err.Error(), "", -1)
	}

	return wrappedResp(vpn.STATUS_OK, vpn.MESSAGE_OK, item.BindIp, int32(item.State))
}

func (service *LiteVpnService) UnRegister(ctx context.Context, req *pb_gen.UnRegisterRequest) (*pb_gen.UnRegisterResponse, error) {

	wrappedResp := func(code, message string, result bool) (resp *pb_gen.UnRegisterResponse, err error) {
		if code != vpn.STATUS_OK {
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
		return wrappedResp(vpn.STATUS_BADREQUEST, "token can't be empty", false)
	}

	vpnMgr := sqlite.VpnMgr{}
	item, err := vpnMgr.QueryByToken(req.Token)
	if item == nil {
		return wrappedResp(vpn.STATUS_OK, err.Error(), false)
	} else if err != nil {
		return wrappedResp(vpn.STATUS_ERR, err.Error(), false)
	}

	result, err := vpnMgr.DeleteById(item.Id)
	if err != nil {
		return wrappedResp(vpn.STATUS_ERR, err.Error(), result)
	}

	service.unRegisterCh <- item.BindIp
	return wrappedResp(vpn.STATUS_OK, vpn.MESSAGE_OK, result)
}

func (service *LiteVpnService) GetRegistedIp(ctx context.Context, req *pb_gen.GetRegistedIpRequest) (*pb_gen.GetRegistedIpResponse, error) {

	wrappedResp := func(code, message, ip string) (resp *pb_gen.GetRegistedIpResponse, err error) {
		if code != vpn.STATUS_OK {
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
		return wrappedResp(vpn.STATUS_BADREQUEST, "token can't be empty", "")
	}

	vpnMgr := sqlite.VpnMgr{}
	item, err := vpnMgr.QueryByToken(req.Token)
	if item == nil {
		return wrappedResp(vpn.STATUS_OK, err.Error(), "")
	} else if err != nil {
		return wrappedResp(vpn.STATUS_ERR, err.Error(), "")
	}

	return wrappedResp(vpn.STATUS_OK, vpn.MESSAGE_OK, item.BindIp)
}

func (service *LiteVpnService) GetToken(ctx context.Context, req *pb_gen.GetTokenRequest) (*pb_gen.GetTokenResponse, error) {

	wrappedResp := func(code, message, token string) (resp *pb_gen.GetTokenResponse, err error) {
		if code != vpn.STATUS_OK {
			err = errors.New(message)
		}
		resp = &pb_gen.GetTokenResponse{
			Message: message,
			Code:    code,
			Token:   token,
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
		State:  vpn.STATE_IDLE,
		BindIp: "",
	})
	if err != nil {
		return wrappedResp(vpn.STATUS_ERR, err.Error(), "")
	}

	return wrappedResp(vpn.STATUS_OK, vpn.MESSAGE_OK, token)
}
