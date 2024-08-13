package dao

import (
	"container/list"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go-pubchem/pkg"
	"reflect"
	"strconv"
)

func getMysqlCursor(host string, port int, username string, passwd string, dbname string) *sqlx.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=latin1", username, passwd, host, port, dbname)
	db := sqlx.MustConnect("mysql", dsn)

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(4)
	return db
}

var MysqlCursor = getMysqlCursor("192.168.2.139", 2881, "luocx@aidb", "ABab12@#", "enotess")

//var MysqlCursor = getMysqlCursor("192.168.2.139", 6306, "luocx", "ABab12@#", "mysql")

func Insert(s string, params ...interface{}) {
	// 注意：这里必须是params... 不然会转成数组类型了
	res, err := MysqlCursor.Exec(s, params...)
	if res != nil {
		pkg.Logger.Error(err, "sql:", s)
	} else {
		pkg.Logger.Info("sql exec :", s)
	}
}

func Delete(s string, params ...interface{}) {
	res, err := MysqlCursor.Exec(s, params...)
	if err != nil {
		pkg.Logger.Error("删除失败：%v\n", err)
	} else {
		pkg.Logger.Info(res.RowsAffected())
	}
}

func Update(s string, params ...interface{}) {
	res, err := MysqlCursor.Exec(s, params...)
	if err != nil {
		pkg.Logger.Error("更新失败：%v\n", err)
	} else {
		pkg.Logger.Info(res.RowsAffected())
	}
}

func Query(s string, i interface{}) (res *list.List) {
	rows, _ := MysqlCursor.Query(s)
	res = list.New()
	t := reflect.TypeOf(i)
	// 获取表的所有列
	cols, _ := rows.Columns()
	// 这里需要琢磨一下
	values := make([]sql.RawBytes, len(cols))
	scanArgs := make([]interface{}, len(values))
	for i := 0; i < len(values); i++ {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		_ = rows.Scan(scanArgs...)
		obj := reflect.New(t).Elem()
		for i, v := range values {
			// FiledByName 是根据字段名获取该字段的值...
			value := obj.FieldByName(cols[i])
			if !value.IsValid() {
				fmt.Printf("字段名为 %s 的字段与结构体不匹配\n", cols[i])
				continue
			} else {
				switch value.Kind() {
				case reflect.Bool:
					// TODO 转换Bool类型
				case reflect.String:
					value.SetString(string(v))
				case reflect.Int:
					temp, _ := strconv.Atoi(string(v))
					value.SetInt(int64(temp))
				}
			}
		}
		res.PushBack(obj)
	}
	return
}

//func BuildInsertQuery(compound pkg.Compound) string {
//	columns := []string{}
//	values := []interface{}{}
//	placeholders := []string{}
//
//	for key, value := range compound {
//		columns = append(columns, key)
//		values = append(values, value)
//		placeholders = append(placeholders, "?")
//	}
//
//	query := fmt.Sprintf("INSERT INTO your_table_name (%s) VALUES (%s)", strings.Join(columns, ","), strings.Join(placeholders, ","))
//
//	return query
//}
