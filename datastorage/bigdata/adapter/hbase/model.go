package hbase

import (
	"fmt"
	. "github.com/arasuresearch/arasu/datastorage/bigdata/adapter/hbase/thrift/Hbase"
	"reflect"
	"strings"
)

type Scope struct {
	Value interface{}
	//Search     *search
	Sql     string
	SqlVars []interface{}
	//db         *DB
	_values    map[string]interface{}
	skipLeft   bool
	primaryKey string
}

func (a *HbaseAdapter) GetRows(v interface{}) error {
	value := reflect.Indirect(reflect.ValueOf(v))
	tname := a.DbName + ":" + value.Type().Elem().Name()
	tscan := NewTScan()
	scanner, err := a.Conn.ScannerOpenWithScan(Text(tname), tscan, nil)
	if err != nil {
		return err
	}
	var cf, cq, Id string
	records := make(map[string]map[string]map[string]interface{})
	for {
		res, err := a.Conn.ScannerGetList(scanner, int32(10))
		if err != nil {
			return err
		}
		if len(res) == 0 {
			break
		}
		for _, e := range res {
			Id = string(e.Row)
			if records[Id] == nil {
				records[Id] = make(map[string]map[string]interface{})
			}
			for k, v := range *e.Columns {
				column := strings.SplitN(k, ":", 2)
				cf = column[0]
				cq = column[1]
				if records[Id][cf] == nil {
					records[Id][cf] = make(map[string]interface{})
				}
				records[Id][cf][cq] = string(v.Value)
			}
		}
	}
	for k, v := range records {
		record := reflect.New(value.Type().Elem()).Elem()
		record.FieldByName("Id").Set(reflect.ValueOf(k))
		for k0, v0 := range v {
			field := record.FieldByName(k0)
			field.Set(reflect.ValueOf(v0))
		}
		value.Set(reflect.Append(value, record))
	}

	if err := a.Conn.ScannerClose(scanner); err != nil {
		return err
	}
	return nil
}
func (a *HbaseAdapter) GetRow(v interface{}) error {
	value := reflect.Indirect(reflect.ValueOf(v))
	tname := a.DbName + ":" + value.Type().Name()
	Id := value.FieldByName("Id").String()
	res, err := a.Conn.GetRow(Text(tname), Text(Id), nil)
	if err != nil {
		return err
	}
	var cf, cq string
	for _, e := range res {
		cfs := make(map[string]map[string]interface{})
		for k, v := range *e.Columns {
			column := strings.SplitN(k, ":", 2)
			cf = column[0]
			cq = column[1]
			if cfs[cf] == nil {
				cfs[cf] = make(map[string]interface{})
			}
			cfs[cf][cq] = string(v.Value)
		}
		for k, v := range cfs {
			field := value.FieldByName(k)
			field.Set(reflect.ValueOf(v))
		}
		value.FieldByName("Id").Set(reflect.ValueOf(string(e.Row)))
	}
	return nil
}
func (a *HbaseAdapter) DeleteRow(v interface{}) error {
	value := reflect.Indirect(reflect.ValueOf(v))
	tname := a.DbName + ":" + value.Type().Name()

	Id := value.FieldByName("Id").String()
	return a.Conn.DeleteAllRow(Text(tname), Text(Id), nil)
}

func (a *HbaseAdapter) Save(v interface{}) error {
	value := reflect.Indirect(reflect.ValueOf(v))
	tname := a.DbName + ":" + value.Type().Name()
	mutations := []*Mutation{}

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		cfName := value.Type().Field(i).Name
		cfs, ok := field.Interface().(map[string]interface{})
		//fmt.Println(cfName, cfs, ok)

		if ok {
			for k, v := range cfs {
				c := cfName + ":" + k
				mutation := Mutation{Column: Text(c), Value: Text(fmt.Sprint(v))}
				mutations = append(mutations, &mutation)
			}
		}
	}
	rowField := value.FieldByName("Id")
	//fmt.Println(rowField, rowField.Interface() == "")
	//fmt.Println(a.GenUUID())

	if rowField.Interface() == "" {
		rowField.Set(reflect.ValueOf(a.GenUUID()))
	}

	Id := value.FieldByName("Id").Interface().(string)
	//fmt.Println(Id)
	return a.Conn.MutateRow(Text(tname), Text(Id), mutations, nil)
}

//record.Field(0).Set(reflect.ValueOf(m1.Basic))

//mms := []Mike{mike}
//ee := reflect.ValueOf(mms)
//r.Elem().Set(ee)

//	GetRows(tableName Text, rows []Text, attributes map[string]Text) (r []*TRowResult_, err error)

// type TRowResult_ struct {
// 	Row           Text               `thrift:"row,1"`
// 	Columns       *map[string]*TCell `thrift:"columns,2"`
// 	SortedColumns *[]*TColumn        `thrift:"sortedColumns,3"`
// }
// type TCell struct {
// 	Value     Bytes `thrift:"value,1"`
// 	Timestamp int64 `thrift:"timestamp,2"`
// }
// var mike Mike
// r := reflect.ValueOf(mike)
// i := reflect.Indirect(r)
// fmt.Println(i.NumField())

// var mike *Mike
// mike = new(Mike)
// r := reflect.ValueOf(mike)
// i := reflect.Indirect(r)
// fmt.Println(i.NumField())

// var mike []Mike
// r := reflect.ValueOf(mike)
// d := reflect.Indirect(r)
// data := reflect.New(d.Type().Elem()).Elem()
// fmt.Println(d.Type(), d.Type().Elem(), data.NumField())

// data := reflect.Indirect(reflect.ValueOf(scope.Value))

// if data.Kind() == reflect.Struct {
// 	if field := data.FieldByName(snakeToUpperCamel(scope.PrimaryKey())); field.IsValid() {
// 		return field.Interface()
// 	}
// }
// return 0
//			data = reflect.New(data.Type().Elem()).Elem()
