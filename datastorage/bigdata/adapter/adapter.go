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
	"github.com/arasuresearch/arasu/datastorage/bigdata/adapter/abstract"
	"github.com/arasuresearch/arasu/datastorage/bigdata/adapter/hbase"
	// "github.com/arasuresearch/arasu/datastorage/bigdata/adapter/bigtable"
	// "github.com/arasuresearch/arasu/datastorage/bigdata/adapter/hypertable"
	// "github.com/arasuresearch/arasu/datastorage/bigdata/adapter/mongodb"
	"log"
)

type Adapter interface {
	GetDbName() string
	CreateDatabase() error
	DropDatabase() error

	GetTableNames() ([]string, error)
	IsDatabaseExists() bool

	CreateTable(string, func(*abstract.Table)) error
	AlterTable(name string, args ...interface{}) error
	DropTable(string) error

	//TableExists(name string) (bool, error)

	// CreateColumn(table_name string, column_name string, args ...interface{})
	// AlterColumn(table_name string, old_column_name string, new_column_name string, args ...interface{})
	// DropColumn(table_name string, column_name string, args ...interface{})
	DumpSchema()

	CreateSchemaMigration() error
	DropSchemaMigration() error
	GetAllSchemaMigration() ([]string, error)
	DeleteFromSchemaMigration(version string) error
	InsertIntoSchemaMigration(version string) error
	SchemaToStruct(string) (*abstract.SchemaToStruct, error)
}

func New(name string, conf string) (adapter Adapter) {
	var abstractAdapter = abstract.New(name, conf)
	switch name {
	case "hbase":
		adapter = hbase.New(abstractAdapter)
	// case "bigtable":
	// 	conn.Adapter = bigtable.New(Conn)
	// case "hypertable":
	// 	conn.Adapter = hypertable.New(Conn)
	// case "mongodb":
	// 	conn.Adapter = mongodb.New()

	default:
		log.Fatal("Adapter Error (" + name + ") Adapter not supported ")
	}
	return
}
