package sqlite

import "errors"

type VpnMgr struct {
	Id     int64
	Token  string
	State  int
	BindIp string
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
		rows.Scan(&id, &token, &state, &bindIp)
		result = append(result, VpnMgr{id, token, state, bindIp})
	}
	if len(result) == 0 {
		return nil, errors.New("fail to find such item")
	}
	return &result[0], nil
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
