// Package src coding=utf-8
// @Project : go-pubchem
// @Time    : 2025/1/13 14:41
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : db_manager.go
// @Software: GoLand
package src

import (
	"database/sql"
	"errors"
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

// AddDefaultValueForVarchar
// @Summary AddDefaultValueForVarchar 设置表的字段为非空，并设置默认值
// @Description insert compound info to db by cid
// @Tags db
// @Accept json
// @Param cid body DbAndTableName true "Cid"
// @Success 200 {string} string "{"msg": "hello wy"}"
// @Failure 400 {string} string "{"msg": "who are you"}"
// @Router /db/addDefaultValue [post]
func AddDefaultValueForVarchar(c *gin.Context) {
	var t DbAndTableName
	err := c.ShouldBind(&t)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}
	type tblInfo struct {
		TableName        string         `db:"table_name"`
		ColumnName       string         `db:"column_name"`
		DataType         string         `db:"data_type"`
		ColumnLen        sql.NullInt64  `db:"column_len"`
		CharacterSetName sql.NullString `db:"character_set_name"`
		CollationName    sql.NullString `db:"collation_name"`
		ColumnComment    sql.NullString `db:"column_comment"`
		TableRowsCount   int64          `db:"table_rows_count"`
	}

	var tbls []tblInfo
	qSql := fmt.Sprintf(`SELECT CONCAT(c.TABLE_SCHEMA, '.', c.TABLE_NAME) as table_name, 
			   c.column_name as column_name, c.DATA_TYPE as data_type,
			   c.CHARACTER_MAXIMUM_LENGTH as column_len, c.CHARACTER_SET_NAME as character_set_name,
			   c.COLLATION_NAME as collation_name, c.COLUMN_COMMENT as column_comment, t.TABLE_ROWS as table_rows_count
			   FROM information_schema.COLUMNS c left join information_schema.tables t on 
				c.TABLE_SCHEMA = t.TABLE_SCHEMA and c.TABLE_NAME = t.TABLE_NAME 
		WHERE c.TABLE_SCHEMA = '%s' and c.TABLE_NAME = '%s' and (c.IS_NULLABLE = 'YES' or c.COLUMN_DEFAULT is null)`, t.DbName, t.TableName)
	err = dao.AidbCursor.Select(&tbls, qSql)
	if err != nil {
		utils.InternalRequestErr(c, errors.New(qSql))
		return
	}

	if len(tbls) == 0 {
		utils.OkRequest(c, "done")
		return
	}

	for _, tbl := range tbls {
		commentStr := ""
		if tbl.ColumnComment.Valid || len(tbl.ColumnComment.String) > 0 {
			commentStr = fmt.Sprintf(`COMMENT '%s'`, tbl.ColumnComment.String)
		}
		// 第一次先处理小于50万的表的varchar
		var alterTableSql string
		switch tbl.DataType {
		case "varchar":
			alterTableSql = fmt.Sprintf("alter /*+ parallel(8) +*/ table %s.%s MODIFY COLUMN `%s` varchar(%d) CHARACTER SET %s COLLATE %s DEFAULT '' NOT NULL "+commentStr, t.DbName, t.TableName, tbl.ColumnName, tbl.ColumnLen.Int64, tbl.CharacterSetName.String, tbl.CollationName.String)
		case "bigint", "tinyint", "decimal", "double", "int", "float", "smallint":
			alterTableSql = fmt.Sprintf("alter /*+ parallel(8) +*/ table %s.%s MODIFY COLUMN `%s` %s DEFAULT 0 NOT NULL "+commentStr, t.DbName, t.TableName, tbl.ColumnName, tbl.DataType)
		case "enum", "longtext", "mediumtext", "text", "tinytext", "char", "blob", "mediumblob", "longblob", "varbinary", "date", "datetime", "bit":
		}

		if alterTableSql == "" {
			utils.OkRequest(c, "don't nedd change")
			return
		}
		_, err = dao.AidbCursor.Exec(alterTableSql)
		if err != nil {
			utils.InternalRequestErr(c, errors.New(err.Error()+" "+alterTableSql))
			return
		}
	}

	utils.OkRequest(c, "done")
	return
}
