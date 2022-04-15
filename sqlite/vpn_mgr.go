package sqlite

import (
	"errors"
	"time"
)

type VpnMgr struct {
	Id         int64
	Token      string
	State      int
	BindIp     string
	CreateTime time.Time
	UpdateTime time.Time
}

func (vpn *VpnMgr) Insert(u VpnMgr) error {
	db = GetDb()
	sql := `insert into vpn_mgr (token, state, bind_ip) values(?,?,?)`
	stmt, err := db.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(u.Token, u.State, u.BindIp)
	return err
}

func (vpn *VpnMgr) InsertToken(token string) error {
	db = GetDb()
	sql := `insert into vpn_mgr (token) values()`
	stmt, err := db.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(token)
	return err
}

func (vpn *VpnMgr) QueryAll() (bindIps []string, e error) {
	db = GetDb()
	sql := `select bind_ip from vpn_mgr`
	stmt, err := db.Prepare(sql)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	var result = make([]string, 0)
	for rows.Next() {
		var bindIp string
		rows.Scan(&bindIp)
		result = append(result, bindIp)
	}
	return result, nil
}

func (vpn *VpnMgr) QueryByToken(token string) (l *VpnMgr, e error) {
	db = GetDb()
	sql := `select * from vpn_mgr where token=?`
	stmt, err := db.Prepare(sql)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(token)
	if err != nil {
		return nil, err
	}
	var result = make([]VpnMgr, 0)
	for rows.Next() {
		var token, bindIp string
		var id int64
		var state int
		var createTime, updateTime time.Time
		rows.Scan(&id, &token, &state, &bindIp, &createTime, &updateTime)
		result = append(result, VpnMgr{id, token, state, bindIp, createTime, updateTime})
	}
	if len(result) == 0 {
		return nil, errors.New("fail to find such item")
	}
	return &result[0], nil
}

func (vpn *VpnMgr) QueryByIp(ip string) (l *VpnMgr, e error) {
	db = GetDb()
	sql := `select * from vpn_mgr where bind_ip=?`
	stmt, err := db.Prepare(sql)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(ip)
	if err != nil {
		return nil, err
	}
	var result = make([]VpnMgr, 0)
	for rows.Next() {
		var token, bindIp string
		var id int64
		var state int
		var createTime, updateTime time.Time
		rows.Scan(&id, &token, &state, &bindIp, &createTime, &updateTime)
		result = append(result, VpnMgr{id, token, state, bindIp, createTime, updateTime})
	}
	if len(result) == 0 {
		return nil, errors.New("fail to find such item")
	}
	return &result[0], nil
}

func (vpn *VpnMgr) QueryLogestIdle() (l *VpnMgr, e error) {
	db = GetDb()
	// valid in sqlite
	sql := `select id,token,state,bind_ip,create_time,update_time,min(update_time) from vpn_mgr where state=-1`

	stmt, err := db.Prepare(sql)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	var result = make([]VpnMgr, 0)
	for rows.Next() {
		var token, bindIp string
		var id int64
		var state int
		var createTime, updateTime, tmp time.Time
		rows.Scan(&id, &token, &state, &bindIp, &createTime, &updateTime, &tmp)
		result = append(result, VpnMgr{id, token, state, bindIp, createTime, updateTime})
		//rows.Scan(&id, &token, &bindIp, &updateTime)
		//result = append(result, VpnMgr{
		//	Id:         id,
		//	Token:      token,
		//	BindIp:     bindIp,
		//	UpdateTime: updateTime,
		//})
	}
	if len(result) == 0 {
		return nil, errors.New("fail to find such item")
	}
	return &result[0], nil
}

func (vpn *VpnMgr) UpdateStateByToken(state int, token string) (bool, error) {
	db = GetDb()
	sql := `update vpn_mgr set state=? where token=?`
	stmt, err := db.Prepare(sql)
	if err != nil {
		return false, err
	}
	res, err := stmt.Exec(state, token)
	if err != nil {
		return false, err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (vpn *VpnMgr) UpdateIpByToken(ip, token string) (bool, error) {
	db = GetDb()
	sql := `update vpn_mgr set bind_ip=? where token=?`
	stmt, err := db.Prepare(sql)
	if err != nil {
		return false, err
	}
	res, err := stmt.Exec(ip, token)
	if err != nil {
		return false, err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (vpn *VpnMgr) DeleteById(id int64) (bool, error) {
	db = GetDb()
	sql := `delete from vpn_mgr where id=?`
	stmt, err := db.Prepare(sql)
	if err != nil {
		return false, err
	}
	res, err := stmt.Exec(id)
	if err != nil {
		return false, err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (vpn *VpnMgr) DeleteByToken(token string) (bool, error) {
	db = GetDb()
	sql := `delete from vpn_mgr where token=?`
	stmt, err := db.Prepare(sql)
	if err != nil {
		return false, err
	}
	res, err := stmt.Exec(token)
	if err != nil {
		return false, err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return false, err
	}
	return true, nil
}
