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
