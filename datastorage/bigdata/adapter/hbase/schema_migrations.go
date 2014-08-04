package hbase

import (
	. "github.com/arasuresearch/arasu/datastorage/bigdata/adapter/hbase/thrift/Hbase"
	"strings"
)

func (a *HbaseAdapter) GetSchemaMigrationConf(args ...string) (t, r, cf, c Text, attrs map[string]Text) {
	t = Text(a.DbName + ":" + a.SchemaTableName)
	cf = Text(a.SchemaTableColumnFamily)
	r = Text(a.SchemaTableVersion)
	if len(args) > 0 {
		c = Text(a.SchemaTableColumnFamily + ":" + args[0])
	}
	return
}

//creating schema_migrations table to take the migration history and order
// arasu@ubuntu:~/projects/demo$ arasu ds bd migrate

func (a *HbaseAdapter) CreateSchemaMigration() error {
	t, _, cf, _, _ := a.GetSchemaMigrationConf()
	cd := NewColumnDescriptor()
	cd.Name = cf
	return a.Conn.CreateTable(t, []*ColumnDescriptor{cd})
}

//droping schema_migrations table to take the migration history and order
// arasu@ubuntu:~/projects/demo$ arasu ds bd rollback
func (a *HbaseAdapter) DropSchemaMigration() error {
	t, _, _, _, _ := a.GetSchemaMigrationConf("")
	if err := a.Conn.DisableTable(Bytes(t)); err != nil {
		return err
	}
	return a.Conn.DeleteTable(t)
}
func (a *HbaseAdapter) InsertIntoSchemaMigration(version string) error {
	t, r, _, c, attrs := a.GetSchemaMigrationConf(version)
	mutations := []*Mutation{&Mutation{Column: c, Value: Text("")}}
	return a.Conn.MutateRow(t, r, mutations, attrs)
}
func (a *HbaseAdapter) DeleteFromSchemaMigration(version string) error {
	t, r, _, c, attrs := a.GetSchemaMigrationConf(version)
	return a.Conn.DeleteAll(t, r, c, attrs)
}

func (a *HbaseAdapter) GetAllSchemaMigration() ([]string, error) {
	ok, err := a.IsThisTableExists(a.SchemaTableName)
	if err != nil {
		return nil, err
	}
	if !ok {
		if err := a.CreateSchemaMigration(); err != nil {
			return nil, err
		}
	}
	t, r, _, _, attrs := a.GetSchemaMigrationConf()
	res, err := a.Conn.GetRow(t, r, attrs)
	if err != nil {
		return nil, err
	}
	result := []string{}
	for _, e := range res {
		for k, _ := range *e.Columns {
			v := strings.TrimLeft(k, a.SchemaTableColumnFamily+":")
			result = append(result, v)
		}
	}
	return result, nil
}

// func (a *AbstractAdapter) GetAllSchemaMigration() ([]string, error) {
// 	ok, err := a.TableExists("schema_migrations")
// 	if err != nil {
// 		return nil, err
// 	}
// 	if !ok {
// 		if err := a.CreateSchemaMigration(); err != nil {
// 			return nil, err
// 		}
// 	}

// 	sql := "SELECT " + a.Quote("schema_migrations") + ".* FROM " + a.Quote("schema_migrations")
// 	return a.QueryRowsFirstColumnStringArray(sql)
// }

func (a *HbaseAdapter) DumpSchema() {

}
