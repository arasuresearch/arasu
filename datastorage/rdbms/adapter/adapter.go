// --
// The BSD License (BSD)

// Copyright (c) 2015 Thaniyarasu Kannusamy <thaniyarasu@gmail.com> & Arasu Research Lab Pvt Ltd. All rights reserved.

// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:

//    * Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above copyright notice, this list of
//    conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
//    * Neither Thaniyarasu Kannusamy <thaniyarasu@gmail.com>. nor ArasuResearch Inc may be used to endorse or promote products derived from this software without specific prior written permission.

// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND AUTHOR
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR AUTHOR BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
// ++

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
