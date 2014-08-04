package hbase

import (
	"fmt"
	"github.com/arasuresearch/arasu/datastorage/bigdata/adapter/abstract"
	. "github.com/arasuresearch/arasu/datastorage/bigdata/adapter/hbase/thrift/Hbase"
	"github.com/arasuresearch/arasu/lib"
	"io/ioutil"
	"path"
	"strings"
)

func (a *HbaseAdapter) SchemaToStruct(adapter_path string) (*abstract.SchemaToStruct, error) {
	model_names, err := ioutil.ReadDir(adapter_path)
	if err != nil {
		return nil, err
	}
	models := lib.AS{}
	for _, e := range model_names {
		if !e.IsDir() && e.Name() != "init.go" && e.Name() != "schema_structures.go" {
			models.Add(strings.TrimSuffix(e.Name(), ".go"))
		}
	}

	imports := lib.AS{}
	tables := []abstract.TableStruct{}
	fmt.Println("updating schema structures for -> ", models)
	for _, e := range models {
		table := abstract.TableStruct{Name: strings.Title(e)}

		column := abstract.ColumnStruct{Name: "Id", Type: "string"}
		table.Columns = append(table.Columns, column)

		tname := Text(a.DbName + ":" + e)
		cds, err := a.Conn.GetColumnDescriptors(tname)
		if err != nil {
			return nil, err
		}
		for k, v := range cds {
			_ = k
			name := string(v.Name)
			name = strings.TrimRight(name, ":")

			column := abstract.ColumnStruct{Name: name, Type: "map[string]interface{}"}
			table.Columns = append(table.Columns, column)
		}
		tables = append(tables, table)
	}
	_, pkg := path.Split(adapter_path)
	data := abstract.SchemaToStruct{pkg, imports, tables}
	return &data, nil
}
