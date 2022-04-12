package sqlite

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"ws-vpn/utils"
)

const (
	dbDriverName = "sqlite3"
	dbName       = "/tmp/litevpn.db"
)

var db *sql.DB
var logger = utils.GetLogger()

func GetDb() *sql.DB {
	if db == nil {
		InitSqlite()
	}
	return db
}

func InitSqlite() (err error) {
	db, err = sql.Open(dbDriverName, dbName)
	if err != nil {
		logger.Info("fail to open sqlite err: %+v", err)
		return
	}
	err = createTable()
	if err != nil {
		logger.Info("fail to create table: %+v", err)
		return
	}
	return
}

func createTable() error {
	sql := `create table if not exists "vpn_mgr" (
		"id" integer primary key autoincrement,
		"token" text not null unique,
		"state" integer not null,
		"bind_ip" text unique
	)`
	_, err := db.Exec(sql)
	return err
}
