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

package abstract

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/url"
	"strconv"
	"strings"
)

type AbstractAdapter struct {
	Name                string
	Conf                string
	ConfWithoutDb       string
	DB                  *sql.DB
	Tx                  *sql.Tx
	TxErr               error
	DbName              string
	Options             map[string][]string
	NativeDatabaseTypes map[string]interface{}
	QuotedColumnNames   map[string]string
	QuotedTableNames    map[string]string
	MaxIndexLength      int
	Credentials         string
}

func (a *AbstractAdapter) Close() error {
	return a.DB.Close()
}
func (a *AbstractAdapter) RaiseIfErrorOccured(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func New(name string, conf string) (abstractAdapter AbstractAdapter) {
	abstractAdapter = AbstractAdapter{
		NativeDatabaseTypes: NATIVE_DATABASE_TYPES,
		QuotedColumnNames:   QUOTED_COLUMN_NAMES,
		QuotedTableNames:    QUOTED_TABLE_NAMES,
		MaxIndexLength:      32,
	}

	confWithoutDb, dbName, options := parseConfiguration(conf)
	abstractAdapter.ConfWithoutDb = confWithoutDb
	abstractAdapter.DbName = dbName
	abstractAdapter.Options = options
	abstractAdapter.Name = name
	abstractAdapter.Conf = conf
	return
}
func (a *AbstractAdapter) InitWithoutDb() error {

	if db, err := sql.Open(a.Name, a.ConfWithoutDb); err != nil {
		return err
	} else {
		a.DB = db
	}
	if err := a.DB.Ping(); err != nil {
		return err
	}
	return nil
}
func (a *AbstractAdapter) Init() error {
	if db, err := sql.Open(a.Name, a.Conf); err != nil {
		return err
	} else {
		a.DB = db
	}
	if err := a.DB.Ping(); err != nil {
		if strings.Contains(err.Error(), "Unknown database") {
			return errors.New(err.Error() + "database (" + a.DbName + ") does not exists .create it first by 'arasu db rdbms create' command")
		}
		return err
	}
	return nil
}
func (a *AbstractAdapter) Transaction(Func func()) error {
	txn, err := a.DB.Begin()
	if err != nil {
		return err
	}
	a.Tx = txn
	a.TxErr = nil
	Func()
	if err := a.TxErr; err != nil {
		a.TxErr = nil
		txn.Rollback()
		return err
	}
	txn.Commit()
	return nil
}
func (a *AbstractAdapter) DumpSchema() error {
	return nil
}

func (a *AbstractAdapter) TypeToSql(ctype string, limit uint64, precision int, scale int) string {
	ctype_native_options := a.NativeDatabaseTypes[ctype].(map[string]interface{})
	column_type_sql := ctype_native_options["name"].(string)

	if ctype == "decimal" {
		if precision > 0 && scale > 0 {
			column_type_sql += "(" + strconv.Itoa(precision) + "," + strconv.Itoa(scale) + ")"
		} else if precision > 0 {
			column_type_sql += "(" + strconv.Itoa(precision) + ")"
		} else if scale > 0 {
			log.Fatal("Error adding decimal column: precision cannot be empty if scale is specified")
		}

	} else if ctype != "primary_key" && limit > 0 {
		column_type_sql += "(" + strconv.Itoa(int(limit)) + ")"
	}
	return column_type_sql
}
func (a *AbstractAdapter) QuoteColumnName(name string) string {
	sql, ok := a.QuotedColumnNames[name]
	if !ok {
		sql = "`" + name + "`"
		a.QuotedColumnNames[name] = sql
	}
	return sql
}
func (a *AbstractAdapter) Quote(name string) string {
	return "`" + name + "`"
}
func (a *AbstractAdapter) SingleQuote(name string) string {
	return "'" + name + "'"
}
func (a *AbstractAdapter) QuoteTableName(name string) string {
	sql, ok := a.QuotedTableNames[name]
	if !ok {
		sql = "`" + name + "`"
		a.QuotedTableNames[name] = sql
	}
	return sql

}

func (a *AbstractAdapter) AutoIncrement() bool {
	return true
}
func (a *AbstractAdapter) SupportsBulkAlter() bool {
	return true
}
func (a *AbstractAdapter) SupportPartialIndex() bool {
	return true
}
func (a *AbstractAdapter) SupportsIndexSortOrder() bool {
	return true
}

var ErrNoDb = errors.New("arasu: no database specified (like 'username:password@protocol(host:post)/database_name?arg=value&...')")

func parseConfiguration(conf string) (conf_without_db string, db_name string, res map[string][]string) {
	u, err := url.Parse(conf)
	if err != nil {
		log.Fatal(err)
	}
	if strings.Contains(u.Opaque, "/") {
		opaques := strings.Split(u.Opaque, "/")
		if len(opaques) > 1 && len(opaques[1]) > 0 {
			conf_without_db = strings.Split(conf, "/")[0] + "/" //"/?" + u.RawQuery //opaques[0]
			db_name = opaques[1]
		} else {
			log.Fatal(ErrNoDb)
		}

	} else {
		log.Fatal(ErrNoDb)
	}

	q := u.Query()
	res = map[string][]string{}
	for k, v := range q {
		res[k] = []string{}
		for _, e := range v {
			if strings.Contains(e, ",") {
				for _, e0 := range strings.Split(e, ",") {
					res[k] = append(res[k], e0)
				}
			} else {
				res[k] = append(res[k], e)
			}

		}
	}
	return
}

var (
	DEFAULT_CHAR_SET      = "utf8"
	QUOTED_COLUMN_NAMES   = map[string]string{}
	QUOTED_TABLE_NAMES    = map[string]string{}
	NATIVE_DATABASE_TYPES = map[string]interface{}{
		"primary_key": map[string]interface{}{"as": "int(11) DEFAULT NULL auto_increment PRIMARY KEY"},
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
)
