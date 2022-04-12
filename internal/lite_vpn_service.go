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
}

var logger = utils.GetLogger()

func CheckConnState(ctx context.Context, req *pb_gen.CheckConnStateRequest) (*pb_gen.CheckConnResponse, error) {

	wrappedResp := func(code, message string, state int32) (resp *pb_gen.CheckConnResponse, err error) {
		if code != vpn.STATUS_OK {
			logger.Errorf("query token: %+v err: %+v", req.Token, err)
			err = errors.New(message)
		}
		resp = &pb_gen.CheckConnResponse{
			Message:   message,
			Code:      code,
			ConnState: state,
		}
		logger.Debugf("resp: %+v", resp)
		return
	}

	if len(req.Token) == 0 {
		return wrappedResp(vpn.STATUS_BADREQUEST, "token can't be empty", -1)
	}

	vpnMgr := sqlite.VpnMgr{}
	item, err := vpnMgr.QueryByToken(req.Token)
	if item == nil {
		return wrappedResp(vpn.STATUS_OK, err.Error(), -1)
	} else if err != nil {
		return wrappedResp(vpn.STATUS_ERR, err.Error(), -1)
	}

	return wrappedResp(vpn.STATUS_OK, vpn.MESSAGE_OK, int32(item.State))
}

func UnRegister(ctx context.Context, req *pb_gen.UnRegisterRequest) (*pb_gen.UnRegisterResponse, error) {

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

	return wrappedResp(vpn.STATUS_OK, vpn.MESSAGE_OK, result)
}
