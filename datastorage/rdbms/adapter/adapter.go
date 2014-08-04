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

package adapter

import (
	"github.com/arasuresearch/arasu/datastorage/rdbms/adapter/abstract"
	"github.com/arasuresearch/arasu/datastorage/rdbms/adapter/mysql"
	// "github.com/arasuresearch/arasu/datastorage/rdbms/adapter/oracle"
	// "github.com/arasuresearch/arasu/datastorage/rdbms/adapter/postgresql"
	// "github.com/arasuresearch/arasu/datastorage/rdbms/adapter/sqlite3"
	"log"
)

type Adapter interface {
	CreateDatabase() error
	DropDatabase() error
	Close() error
	CreateTable(name string, args ...interface{}) error
	AlterTable(name string, args ...interface{}) error
	DropTable(name string, args ...interface{}) error
	TableExists(name string) (bool, error)
	InitWithoutDb() error
	Init() error
	Transaction(func()) error
	DumpSchema() error

	// CreateColumn(table_name string, column_name string, args ...interface{})
	// AlterColumn(table_name string, old_column_name string, new_column_name string, args ...interface{})
	// DropColumn(table_name string, column_name string, args ...interface{})

	CreateIndex(table_name string, column_names []string, args map[string]interface{}) error
	//AlterIndex(table_name string, columns interface{}, args ...interface{}) error
	DropIndex(table_name string, column_names []string, args map[string]interface{}) error

	//CreateSchemaMigration() error
	//DropSchemaMigration() error
	GetAllSchemaMigration() ([]string, error)
	DeleteFromSchemaMigration(version string) error
	InsertIntoSchemaMigration(version string) error
	SchemaToStruct(string) (*abstract.SchemaToStruct, error)
}

func New(name string, conf string) (adapter Adapter) {
	var abstractAdapter = abstract.New(name, conf)
	switch name {
	case "mysql":
		adapter = mysql.New(abstractAdapter)
	// case "postgresql":
	// 	adapter = postgresql.New(abstractAdapter)
	// case "sqlite3":
	// 	adapter = sqlite3.New(abstractAdapter)
	// case "oracle":
	// 	adapter = oracle.New(abstractAdapter)
	default:
		log.Fatal("Adapter Error (" + name + ") Adapter not supported ")
	}
	return
}

// func InitMigrater(name string, conf string) {
// 	_, db_name, query := parse_configuration(conf)

// 	Conn = &abstract.Connection{Name: name, Conf: conf, DbName: db_name, Query: query}

// 	if db, err := sql.Open(name, conf); err == nil {
// 		Conn.DB = db
// 	} else {
// 		log.Fatal(err)
// 	}
// 	if err := Conn.DB.Ping(); err != nil {
// 		if strings.Contains(err.Error(), "Unknown database") {
// 			log.Fatal("First create database (" + db_name + ") by 'arasu db rdbms create'")
// 		} else {
// 			log.Fatal(err)
// 		}
// 		return
// 	}
// 	InitAdapter(name)
// 	if exists, _ := Conn.Adapter.TableExists("schema_migrations"); !exists {
// 		_ = Conn.Adapter.CreateSchemaMigration()
// 	}
// }

// func InitAdapter(name string) {
// 	switch name {
// 	case "mysql":
// 		Conn.Adapter = mysql.New(Conn)
// 	// case "postgresql":
// 	// 	conn.Adapter = postgresql.New()
// 	// case "sqlite3":
// 	// 	conn.Adapter = sqlite3.New()
// 	// case "oracle":
// 	// 	conn.Adapter = oracle.New()
// 	default:
// 		log.Fatal("Connection Error (" + name + ") Adapter not supported ")
// 	}
// }

// // var (
// // 	CreateTable = Conn.Adapter.CreateTable
// // 	//CreateSchemaVersion = Conn.Adapter.CreateSchemaVersion
// // )

// func CreateDatabase(name string, conf string) {
// 	conf_without_db, db_name, query := parse_configuration(conf)
// 	fmt.Println(conf_without_db)

// 	Conn = &abstract.Connection{Name: name, Conf: conf, DbName: db_name, Query: query}

// 	if db, err := sql.Open(name, conf_without_db); err == nil {
// 		Conn.DB = db
// 	} else {
// 		log.Fatal(err)
// 	}
// 	if err := Conn.DB.Ping(); err != nil {
// 		log.Fatal(err)
// 		return
// 	}
// 	InitAdapter(name)
// 	Conn.Adapter.CreateDatabase()
// }
// func DropDatabase(name string, conf string) {

// 	conf_without_db, db_name, query := parse_configuration(conf)
// 	fmt.Println(conf_without_db)

// 	Conn = &abstract.Connection{Name: name, Conf: conf, DbName: db_name, Query: query}

// 	if db, err := sql.Open(name, conf_without_db); err == nil {
// 		Conn.DB = db
// 	} else {
// 		log.Fatal(err)
// 	}
// 	if err := Conn.DB.Ping(); err != nil {
// 		log.Fatal(err)
// 		return
// 	}
// 	InitAdapter(name)
// 	Conn.Adapter.DropDatabase()
// }
