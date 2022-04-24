package sqlite

import (
	"testing"
)

func TestInit(t *testing.T) {
	InitSqlite()
	network := &NetworkMgr{}
	err := network.Insert(NetworkMgr{
		Token:  "999token",
		State:  -1,
		BindIp: "10.1.1.3",
	})
	logger.Info(err)
}

func TestNetworkMgr_QueryByToken(t *testing.T) {
	InitSqlite()
	network := &NetworkMgr{}
	networks, err := network.QueryByToken("wannatoken")
	logger.Infof("%+v", networks)
	logger.Info(err)
}

func TestNetworkMgr_QueryByIp(t *testing.T) {
	InitSqlite()
	network := &NetworkMgr{}
	networks, err := network.QueryByIp("10.1.1.3")
	logger.Infof("%+v", networks)
	logger.Info(err)
}

func TestNetworkMgr_UpdateByIp(t *testing.T) {
	InitSqlite()
	network := &NetworkMgr{}
	res, err := network.UpdateStateByToken(11, "wannatoken")
	logger.Infof("%+v", res)
	logger.Info(err)
}

func TestNetworkMgr_QueryLogestIdle(t *testing.T) {
	InitSqlite()
	network := &NetworkMgr{}
	res, err := network.QueryLogestIdle()
	logger.Infof("%+v", res)
	logger.Info(err)
}
