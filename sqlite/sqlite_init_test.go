package sqlite

import (
	"testing"
)

func TestInit(t *testing.T) {
	InitSqlite()
	vpn := &VpnMgr{}
	err := vpn.Insert(VpnMgr{
		Token:  "999token",
		State:  -1,
		BindIp: "10.1.1.3",
	})
	logger.Info(err)
}

func TestVpnMgr_QueryByToken(t *testing.T) {
	InitSqlite()
	vpn := &VpnMgr{}
	vpns, err := vpn.QueryByToken("wannatoken")
	logger.Infof("%+v", vpns)
	logger.Info(err)
}

func TestVpnMgr_QueryByIp(t *testing.T) {
	InitSqlite()
	vpn := &VpnMgr{}
	vpns, err := vpn.QueryByIp("10.1.1.3")
	logger.Infof("%+v", vpns)
	logger.Info(err)
}

func TestVpnMgr_UpdateByIp(t *testing.T) {
	InitSqlite()
	vpn := &VpnMgr{}
	res, err := vpn.UpdateStateByToken(11, "wannatoken")
	logger.Infof("%+v", res)
	logger.Info(err)
}

func TestVpnMgr_QueryLogestIdle(t *testing.T) {
	InitSqlite()
	vpn := &VpnMgr{}
	res, err := vpn.QueryLogestIdle()
	logger.Infof("%+v", res)
	logger.Info(err)
}
