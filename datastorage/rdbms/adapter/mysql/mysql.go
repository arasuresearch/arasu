// --
// Copyright (c) 2014 Thaniyarasu Kannusamy <thaniyarasu@gmail.com>.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
// ++
//

package mysql

import (
	"github.com/arasuresearch/arasu/datastorage/rdbms/adapter/abstract"
	"log"
	"strconv"
)

type MysqlAdapter struct {
	abstract.AbstractAdapter
}

func New(abstractAdapter abstract.AbstractAdapter) *MysqlAdapter {
	abstractAdapter.NativeDatabaseTypes = NATIVE_DATABASE_TYPES
	var mysqlAdapter = &MysqlAdapter{AbstractAdapter: abstractAdapter}
	return mysqlAdapter
}

func (m *MysqlAdapter) TypeToSql(ctype string, limit uint64, precision int, scale int) string {
	var sql string
	limit_string := strconv.Itoa(int(limit))

	switch ctype {
	case "binary":
		switch {
		case limit <= 0xfff:
			sql = "varbinary(" + limit_string + ")"
		case limit == 0:
			sql = "blob"
		case limit >= 0x1000 && limit <= 0xffffffff:
			sql = "blob(" + limit_string + ")"
		default:
			log.Fatal("No binary type has character length " + limit_string)
		}

	case "integer":
		switch limit {
		case 1:
			sql = "tinyint"
		case 2:
			sql = "smallint"
		case 3:
			sql = "mediumint"
		case 0, 4, 11:
			sql = "int(11)"
		case 5, 6, 7, 8:
			sql = "bigint"
		default:
			log.Fatal("No integer type has byte size " + limit_string)
		}

	case "text":
		switch {
		case limit > 0 && limit <= 0xff:
			sql = "tinytext"
		case limit == 0 || (limit >= 0x100 && limit <= 0xffff):
			sql = "text"
		case limit >= 0x10000 && limit <= 0xffffff:
			sql = "mediumtext"
		case limit >= 0x1000000 && limit <= 0xffffffff:
			sql = "longtext"
		default:
			log.Fatal("No text type has character length " + limit_string)
		}

	default:
		sql = m.TypeToSql(ctype, limit, precision, scale)
	}
	return sql

}

func (m *MysqlAdapter) VisitAddColumn(c abstract.Column) string {
	return m.AddColumnPosition(m.VisitAddColumn(c), c.Options)
}

func (m *MysqlAdapter) VisitAlterColumn(chg_column abstract.AlterColumn) string {
	column := chg_column.Column
	options := chg_column.Options

	sql_type := m.TypeToSql(chg_column.Ctype, options["limit"].(uint64), options["precision"].(int), options["scale"].(int))
	change_column_sql := "CHANGE " + m.QuoteColumnName(column.Name) + " " + m.QuoteColumnName(options["name"].(string)) + " " + sql_type
	m.AddColumnOptions(change_column_sql, column)
	return m.AddColumnPosition(change_column_sql, options)

}

func (m *MysqlAdapter) AddColumnPosition(sql string, options map[string]interface{}) string {
	if first, exists := options["first"]; exists {
		if first.(bool) {
			sql += " FIRST"
		}
	} else if after, exists := options["after"]; exists {
		sql += " AFTER " + m.QuoteColumnName(after.(string))
	}
	return sql
}

var NATIVE_DATABASE_TYPES = map[string]interface{}{
	"primary_key": map[string]interface{}{"as": "int(11) auto_increment PRIMARY KEY"},
	"string":      map[string]interface{}{"name": "varchar", "limit": uint64(255)},
	"text":        map[string]interface{}{"name": "text"},
	"integer":     map[string]interface{}{"name": "int", "limit": uint64(4)},
	"float":       map[string]interface{}{"name": "float"},
	"decimal":     map[string]interface{}{"name": "decimal"},
	"datetime":    map[string]interface{}{"name": "datetime"},
	"timestamp":   map[string]interface{}{"name": "datetime"},
	"time":        map[string]interface{}{"name": "time"},
	"date":        map[string]interface{}{"name": "date"},
	"binary":      map[string]interface{}{"name": "blob"},
	"boolean":     map[string]interface{}{"name": "tinyint", "limit": uint64(1)},
}
