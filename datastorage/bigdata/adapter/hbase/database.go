package hbase

import (
	"fmt"
	. "github.com/arasuresearch/arasu/datastorage/bigdata/adapter/hbase/thrift/Hbase"
	"log"
)

// creating database aka namespace in hbase shell
// all the tables related to this migrator will be created and maintained inside this database (or namespace)
// the example command
// arasu@ubuntu:~/projects/demo$ arasu ds bd create
//
// ds - means datastorage it can be rdbms,bigdata,localstorage,etc . arasu defaultly supports for rdbms,bigdata

func (a *HbaseAdapter) CreateDatabase() error {
	if a.IsDatabaseExists() {
		return fmt.Errorf("'%s' database(namespace in hbase) already exixts", a.DbName)
	} else {
		fmt.Println("please create a namespace in hbase shell by executing > create_namespace '" + a.DbName + "' . if you want to open bigdata shell then type >hbase shell")
	}
	return nil
	// name, options := one_string_one_options(args...)
	// if len(name) == 0 {
	// 	name = a.Conn.DbName
	// }
	// charset := "utf8"

	// if flag, ok := options["charset"]; ok {
	// 	charset = flag.(string)
	// }
	// sql := "CREATE DATABASE " + a.Quote(name) + " DEFAULT CHARACTER SET " + a.Quote(charset)
	// // if flag, ok := options["collation"]; ok {
	// // 	sql += " COLLATE " + a.Quote(flag.(string))
	// // }
	// if _, err := a.ExecWithTxn(sql); err != nil {
	// 	return err
	// }
	// return nil
}

// droping database or namespace in hbase shell
// arasu@ubuntu:~/projects/demo$ arasu ds bd drop
//
func (a *HbaseAdapter) DropDatabase() error {
	names, err := a.GetTableNames()
	if err != nil {
		return err
	}
	if len(names) > 0 {

		for _, e := range names {
			tableName := a.DbName + ":" + e
			if err := a.Conn.DisableTable(Bytes(tableName)); err != nil {
				return err
			}
			if err := a.Conn.DeleteTable(Text(tableName)); err != nil {
				return err
			}
		}
	}
	fmt.Println("please delete the namespace in hbase shell by executing > drop_namespace '" + a.DbName + "'")
	return nil
	// name, _ := one_string_one_options(args...)
	// if len(name) == 0 {
	// 	name = a.Conn.DbName
	// }
	// sql := "DROP DATABASE IF EXISTS " + a.Quote(name)
	// if _, err := a.ExecWithTxn(sql); err != nil {
	// 	return err
	// }
	// return nil
}

// look wether database or namepsace is exist or not
func (a *HbaseAdapter) IsDatabaseExists() bool {
	ts := NewTScan()
	ts.Columns = &[]Text{Text("info:d")}
	tableName := Text(a.DbName)
	ts.StartRow = &tableName
	ts.StopRow = &tableName
	scanId, err := a.Conn.ScannerOpenWithScan(Text("hbase:namespace"), ts, nil)
	if err != nil {
		log.Fatal(err)
	}
	res, err := a.Conn.ScannerGetList(scanId, int32(10))
	if err != nil {
		log.Fatal(err)
	}
	if len(res) == 0 {
		return false
	}
	return true
}
