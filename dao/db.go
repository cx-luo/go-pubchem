package dao

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func getMysqlCursor(host string, port int, username string, passwd string, dbname string) *sqlx.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=latin1", username, passwd, host, port, dbname)
	db := sqlx.MustConnect("mysql", dsn)

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(4)
	return db
}

var AidbCursor = getMysqlCursor("192.168.2.139", 2881, "luocx@aidb", "ABab12@#", "enotess")
var ElnCursor = getMysqlCursor("192.168.2.139", 2881, "luocx@chemzero", "ABab12@#", "elabx")
