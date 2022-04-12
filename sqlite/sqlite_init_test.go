package sqlite

import (
	"testing"
)

func TestInit(t *testing.T) {
	InitSqlite()
	vpn := &VpnMgr{}
	err := vpn.Insert(VpnMgr{
		Token:  "wannatoken",
		State:  3,
		BindIp: "10.1.1.3",
	})
	logger.Info(err)
}

func TestVpnMgr_Query(t *testing.T) {
	InitSqlite()
	vpn := &VpnMgr{}
	vpns, err := vpn.QueryByToken("wannatoken1")
	logger.Info(vpns)
	logger.Info(err)
}
