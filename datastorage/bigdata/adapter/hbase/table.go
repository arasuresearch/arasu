package hbase

import (
	"fmt"
	"github.com/arasuresearch/arasu/datastorage/bigdata/adapter/abstract"
	. "github.com/arasuresearch/arasu/datastorage/bigdata/adapter/hbase/thrift/Hbase"

	"github.com/arasuresearch/arasu/lib/stringer"
	"log"
	"strings"
)

//creating table by specified by the migration file
// arasu@ubuntu:~/projects/demo$ arasu ds bd migrate

func (a *HbaseAdapter) CreateTable(name string, callback func(t *abstract.Table)) error {
	log.Printf("creating '%s' table...", name)
	tableName := Text(a.DbName + ":" + name)
	table := abstract.Table{}
	callback(&table)
	cds := make([]*ColumnDescriptor, len(table.Columns))
	for i, e := range table.Columns {
		cds[i] = NewColumnDescriptor()
		cds[i].Name = Text(e)
	}
	err := a.Conn.CreateTable(tableName, cds)

	if err != nil {
		return err
	}
	return nil
}

//droping  table by specified by the migration file
// arasu@ubuntu:~/projects/demo$ arasu ds bd rollback

func (a *HbaseAdapter) DropTable(name string) error {
	log.Printf("dropping '%s' table...", name)

	tableName := Text(a.DbName + ":" + name)
	err := a.Conn.DisableTable(Bytes(tableName))
	if err != nil {
		return err
	}
	err = a.Conn.DeleteTable(tableName)
	if err != nil {
		return err
	}

	return nil
}

//alter table by specified by the migration file
// arasu@ubuntu:~/projects/demo$ arasu ds bd migrate
// arasu@ubuntu:~/projects/demo$ arasu ds bd rollback

func (a *HbaseAdapter) AlterTable(name string, args ...interface{}) error {

	var callback func(t *abstract.AlterTable)
	var options = map[string]interface{}{}

	if len(args) == 1 {
		callback = args[0].(func(t *abstract.AlterTable))
	} else if len(args) > 1 {
		options = args[0].(map[string]interface{})
		callback = args[1].(func(t *abstract.AlterTable))
	}
	_, _ = options, callback
	return nil
}

// listing all talbe inside the database or namespace

func (a *HbaseAdapter) GetTableNames() ([]string, error) {
	res, err := a.Conn.GetTableNames()
	if err != nil {
		return nil, err
	}
	names := []string{}
	prefix := a.DbName + ":"
	for _, e := range res {
		name := string(e)
		if strings.HasPrefix(name, prefix) {
			names = append(names, strings.TrimPrefix(name, prefix))
		}
	}
	fmt.Println(names)
	return names, nil
}

func (a *HbaseAdapter) IsThisTableExists(name string) (bool, error) {
	names, err := a.GetTableNames()
	if err != nil {
		return false, err
	}
	return stringer.Contains(names, name), nil
}
