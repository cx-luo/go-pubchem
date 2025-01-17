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
	"go-pubchem/utils"
)

type DbAndTableName struct {
	DbName    string `json:"dbName"`
	TableName string `json:"tableName"`
}

func (dn DbAndTableName) FillNullValues() error {
	type tblInfo struct {
		ColumnName string `db:"column_name"`
		DataType   string `db:"data_type"`
		ColumnLen  int    `db:"column_len"`
	}

	var tbInfos []tblInfo
	err := dao.AidbCursor.Select(&tbInfos, `SELECT  c.column_name as column_name, 
			   ifnull(c.CHARACTER_MAXIMUM_LENGTH , 0) as column_len,
			   c.DATA_TYPE as data_type
			   FROM 
			information_schema.COLUMNS c left join information_schema.tables t on 
				c.TABLE_SCHEMA =t.TABLE_SCHEMA and c.TABLE_NAME =t.TABLE_NAME 
		WHERE c.TABLE_SCHEMA = ? and c.TABLE_NAME = ? and (c.IS_NULLABLE = 'YES' or c.COLUMN_DEFAULT is null)`, dn.DbName, dn.TableName)
	if err != nil {
		return err
	}
	for _, tbInfo := range tbInfos {
		var updateNullToBullstrSql string
		switch tbInfo.DataType {
		case "varchar":
			updateNullToBullstrSql = fmt.Sprintf(`update %s.%s set %s = '' where %s is null`, dn.DbName, dn.TableName, tbInfo.ColumnName, tbInfo.ColumnName)
			break
		case
			"bigint", "tinyint", "decimal", "double", "int", "float", "smallint":
			updateNullToBullstrSql = fmt.Sprintf(`update %s.%s set %s = 0 where %s is null`, dn.DbName, dn.TableName, tbInfo.ColumnName, tbInfo.ColumnName)
			break
		case "enum", "longtext", "mediumtext", "text", "tinytext", "char", "blob", "mediumblob", "longblob", "varbinary", "date", "datetime", "bit":
			// do nothing
			break
		}
		_, err := dao.AidbCursor.Exec(updateNullToBullstrSql)
		if err != nil {
			return err
		}
	}
	return nil
}

func (dn DbAndTableName) SetDefaultValues() error {
	type tblInfo struct {
		ColumnName       string         `db:"column_name"`
		DataType         string         `db:"data_type"`
		CharacterSetName sql.NullString `db:"character_set_name"`
		CollationName    sql.NullString `db:"collation_name"`
		ColumnComment    sql.NullString `db:"column_comment"`
		TableRowsCount   int64          `db:"table_rows_count"`
		ColumnType       string         `db:"column_type"`
	}

	var tbls []tblInfo
	qSql := fmt.Sprintf(`SELECT  c.column_name as column_name, c.DATA_TYPE as data_type,
			   c.CHARACTER_SET_NAME as character_set_name, c.COLUMN_TYPE as column_type,
			   c.COLLATION_NAME as collation_name, c.COLUMN_COMMENT as column_comment, t.TABLE_ROWS as table_rows_count
			   FROM information_schema.COLUMNS c left join information_schema.tables t on 
				c.TABLE_SCHEMA = t.TABLE_SCHEMA and c.TABLE_NAME = t.TABLE_NAME 
		WHERE c.TABLE_SCHEMA = '%s' and c.TABLE_NAME = '%s' and (c.IS_NULLABLE = 'YES' or c.COLUMN_DEFAULT is null)`, dn.DbName, dn.TableName)
	err := dao.AidbCursor.Select(&tbls, qSql)
	if err != nil {
		return err
	}

	if len(tbls) == 0 {
		return nil
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
			alterTableSql = fmt.Sprintf("alter /*+ parallel(8) +*/ table %s.%s MODIFY COLUMN `%s` %s CHARACTER SET %s COLLATE %s DEFAULT '' NOT NULL "+commentStr, dn.DbName, dn.TableName, tbl.ColumnName, tbl.ColumnType, tbl.CharacterSetName.String, tbl.CollationName.String)
			break
		case "bigint", "tinyint", "decimal", "double", "int", "float", "smallint":
			alterTableSql = fmt.Sprintf("alter /*+ parallel(8) +*/ table %s.%s MODIFY COLUMN `%s` %s DEFAULT 0 NOT NULL "+commentStr, dn.DbName, dn.TableName, tbl.ColumnName, tbl.DataType)
			break
		case "enum", "longtext", "mediumtext", "text", "tinytext", "char", "blob", "mediumblob", "longblob", "varbinary", "date", "datetime", "bit":
			break
		}

		if alterTableSql == "" {
			// do nothing
		} else {
			_, err = dao.AidbCursor.Exec(alterTableSql)
			if err != nil {
				return errors.New(err.Error() + " " + alterTableSql)
			}
		}
	}
	return nil
}

func (dn DbAndTableName) ConvertTextOrBit() error {
	type tblInfo struct {
		ColumnName       string         `db:"column_name"`
		DataType         string         `db:"data_type"`
		CharacterSetName sql.NullString `db:"character_set_name"`
		CollationName    sql.NullString `db:"collation_name"`
		ColumnComment    sql.NullString `db:"column_comment"`
		TableRowsCount   int64          `db:"table_rows_count"`
		ColumnType       string         `db:"column_type"`
	}
	var t []tblInfo
	err := dao.AidbCursor.Select(&t, `SELECT  c.column_name as column_name, c.DATA_TYPE as data_type,
			   c.CHARACTER_SET_NAME as character_set_name, c.COLUMN_TYPE as column_type,
			   c.COLLATION_NAME as collation_name, c.COLUMN_COMMENT as column_comment, t.TABLE_ROWS as table_rows_count
			   FROM information_schema.COLUMNS c left join information_schema.tables t on 
				c.TABLE_SCHEMA = t.TABLE_SCHEMA and c.TABLE_NAME = t.TABLE_NAME
		WHERE c.TABLE_SCHEMA = ? and c.TABLE_NAME = ? and c.DATA_TYPE in ('longtext', 'mediumtext', 'text', 'tinytext', 'bit')`, dn.DbName, dn.TableName)
	if err != nil {
		return err
	}

	for _, tbl := range t {
		switch tbl.DataType {
		case "longtext", "mediumtext", "text", "tinytext":
			var maxLen, setLen int64
			qsql := fmt.Sprintf("select max(length(%s)) from %s.%s", tbl.ColumnName, dn.DbName, dn.TableName)
			err := dao.AidbCursor.Get(&maxLen, qsql)
			if err != nil {
				return err
			}
			switch {
			case maxLen <= 60:
				setLen = 63
				break
			case maxLen <= 120:
				setLen = 126
				break
			case maxLen <= 252:
				setLen = 254
				break
			case maxLen <= 510:
				setLen = 512
				break
			case maxLen <= 1000:
				setLen = 1022
				break
			case maxLen <= 2000:
				setLen = 2040
				break
			case maxLen <= 3000:
				setLen = 3070
				break
			case maxLen <= 10000:
				setLen = maxLen + 512
				break
			case maxLen > 10000:
				return nil
			}
			alterSql := fmt.Sprintf("ALTER /*+ parallel(8) +*/ TABLE %s.%s MODIFY COLUMN %s varchar(%d) CHARACTER SET %s COLLATE %s NULL", dn.DbName, dn.TableName, tbl.ColumnName, setLen, tbl.CharacterSetName.String, tbl.CollationName.String)
			_, err = dao.AidbCursor.Exec(alterSql)
			if err != nil {
				return err
			}
		case "bit":
			alterSql := fmt.Sprintf("ALTER TABLE %s.%s MODIFY COLUMN %s tinyint NULL", dn.DbName, dn.TableName, tbl.ColumnName)
			_, err = dao.AidbCursor.Exec(alterSql)
			if err != nil {
				return err
			}
		default:
			return nil
		}
	}
	return nil
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

	err = t.FillNullValues()
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}

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

	err = t.SetDefaultValues()
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	utils.OkRequest(c, "done")
	return
}

// ConvertTextToVarchar
// @Summary ConvertTextToVarchar convert text/bit to varchar/tinyint
// @Description insert compound info to db by cid
// @Tags db
// @Accept json
// @Param cid body DbAndTableName true "Cid"
// @Success 200 {string} string "{"msg": "hello wy"}"
// @Failure 400 {string} string "{"msg": "who are you"}"
// @Router /db/convertTextToVarchar [post]
func ConvertTextToVarchar(c *gin.Context) {
	var t DbAndTableName
	err := c.ShouldBind(&t)
	if err != nil {
		utils.BadRequestErr(c, err)
		return
	}
	err = t.ConvertTextOrBit()
	if err != nil {
		utils.InternalRequestErr(c, err)
		return
	}

	utils.OkRequest(c, "done")
	return
}
