package abstract

import (
	. "github.com/arasuresearch/arasu/datastorage/bigdata/adapter/hbase/thrift/Hbase"
)

type Table struct {
	Name    string
	Columns []string
	Config  map[string]interface{}
}
type AlterTable struct {
	Name string

	AddColumnFamilies   []*ColumnDescriptor
	DropColumnFamilies  []*ColumnDescriptor
	AlterColumnFamilies []*ColumnDescriptor
	Config              map[string]interface{}
}
type DropTable struct {
	Name   Text
	Config map[string]interface{}
}

type AlterColumn struct {
	TableName string
	Name      string
	Config    map[string]interface{}
}
type DropColumn struct {
	TableName string
	Name      string
	Config    map[string]interface{}
}

func (t *Table) ColumnFamily(name string) {
	t.Columns = append(t.Columns, name)
}
