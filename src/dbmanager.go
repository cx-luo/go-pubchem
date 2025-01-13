// Package src coding=utf-8
// @Project : go-pubchem
// @Time    : 2025/1/13 14:41
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : dbmanager.go
// @Software: GoLand
package src

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-pubchem/dao"
	"go-pubchem/pkg"
	"go-pubchem/utils"
)

type DbAndTableName struct {
	DbName    string `json:"dbName"`
	TableName string `json:"tableName"`
}

// FillNullValue
// @Summary FillNullValue 填充表的 null 值
// @Description insert compound info to db by cid
// @Tags db
// @Accept json
// @Param cid body DbAndTableName true "Cid"
// @Success 200 {string} string "{"msg": "hello wy"}"
// @Failure 400 {string} string "{"msg": "who are you"}"
// @Router /db/fillNullValue [post]
func FillNullValue(c *gin.Context) {
	var t DbAndTableName
	err := c.ShouldBind(&t)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}

	type tblInfo struct {
		TableName  string `db:"table_name"`
		ColumnName string `db:"column_name"`
		DataType   string `db:"data_type"`
		ColumnLen  int    `db:"column_len"`
	}

	var tbInfos []tblInfo
	err = dao.AidbCursor.Select(&tbInfos, `SELECT CONCAT(c.TABLE_SCHEMA, '.', c.TABLE_NAME) as table_name, 
			   c.column_name as column_name, 
			   ifnull(c.CHARACTER_MAXIMUM_LENGTH , 0) as column_len,
			   c.DATA_TYPE as data_type
			   FROM 
			information_schema.COLUMNS c left join information_schema.tables t on 
				c.TABLE_SCHEMA =t.TABLE_SCHEMA and c.TABLE_NAME =t.TABLE_NAME 
		WHERE c.TABLE_SCHEMA = ? and c.TABLE_NAME = ? and (c.IS_NULLABLE = 'YES' or c.COLUMN_DEFAULT is null)`, t.DbName, t.TableName)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}
	for _, tbInfo := range tbInfos {
		sema.Acquire(1)
		go func(tb tblInfo) {
			var updateNullToBullstrSql string
			defer sema.Release()
			switch tb.DataType {
			case "varchar":
				updateNullToBullstrSql = fmt.Sprintf(`update %s set %s = '' where %s is null`, tb.TableName, tb.ColumnName, tb.ColumnName)
			case
				"bigint", "tinyint", "decimal", "double", "int", "float", "smallint":
				updateNullToBullstrSql = fmt.Sprintf(`update %s set %s = 0 where %s is null`, tb.TableName, tb.ColumnName, tb.ColumnName)
			case "enum", "longtext", "mediumtext", "text", "tinytext", "char", "blob", "mediumblob", "longblob", "varbinary", "date", "datetime", "bit":
				return
			}
			_, err := dao.AidbCursor.Exec(updateNullToBullstrSql)
			if err != nil {
				pkg.Logger.Error(err.Error()+": %s", updateNullToBullstrSql)
			}
		}(tbInfo)
	}

	sema.Wait()
	utils.OkRequest(c, "done")
	return
}
