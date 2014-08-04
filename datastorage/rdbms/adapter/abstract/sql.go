package abstract

import (
	"database/sql"
	"reflect"
	"strconv"
	"time"
)

func (a *AbstractAdapter) ExecWithTxn(query string) (result sql.Result, err error) {
	txn, err := a.DB.Begin()
	if err != nil {
		return
	}
	result, err = txn.Exec(query)
	if err != nil {
		txn.Rollback()
		return
	}
	txn.Commit()
	return
}

func one_string_one_options(args ...interface{}) (first string, options map[string]interface{}) {
	length := len(args)
	if length == 1 {
		first = args[0].(string)
	} else if length > 1 {
		first = args[0].(string)
		options = args[1].(map[string]interface{})
	}
	return
}
func two_string_one_options(args ...interface{}) (first string, second string, options map[string]interface{}) {
	length := len(args)
	if length == 1 {
		first = args[0].(string)
	} else if length == 2 {
		first = args[0].(string)
		second = args[1].(string)
	} else if length > 2 {
		first = args[0].(string)
		second = args[1].(string)
		options = args[2].(map[string]interface{})
	}
	return
}

func (a AbstractAdapter) QueryRowsFirstColumnStringArray(query string) ([]string, error) {
	records, err := a.Query(query)
	if err != nil {
		return nil, err
	}
	result := []string{}
	for _, e := range records {
		for _, v := range e {
			result = append(result, v)
			break
		}
	}
	return result, nil
}

func (a AbstractAdapter) Query(query string) (records []map[string]string, err error) {
	res, err := a.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	columns, err := res.Columns()
	if err != nil {
		return nil, err
	}
	cl := len(columns)
	for res.Next() {
		record := map[string]string{}
		scanner := make([]interface{}, cl)
		for i := 0; i < cl; i++ {
			var lv interface{}
			scanner[i] = &lv
		}
		if err := res.Scan(scanner...); err != nil {
			return nil, err
		}
		for i, e := range columns {
			rr := reflect.ValueOf(scanner[i])
			rawValue := reflect.Indirect(rr)
			//rawValue := rr
			//if row is null then ignore
			if rawValue.Interface() == nil {
				continue
			}

			rtype := reflect.TypeOf(rawValue.Interface())
			rvalue := reflect.ValueOf(rawValue.Interface())
			switch rtype.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				record[e] = strconv.FormatInt(rvalue.Int(), 10)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				record[e] = strconv.FormatUint(rvalue.Uint(), 10)
			case reflect.Float32, reflect.Float64:
				record[e] = strconv.FormatFloat(rvalue.Float(), 'f', -1, 64)
			case reflect.Slice:
				if rtype.Elem().Kind() == reflect.Uint8 {
					record[e] = string(rawValue.Interface().([]byte))
					break
				}
			case reflect.String:
				record[e] = rvalue.String()
			case reflect.Struct:
				record[e] = rawValue.Interface().(time.Time).Format("2006-01-02 15:04:05.000 -0700")
			case reflect.Bool:
				if rvalue.Bool() {
					record[e] = "1"
				} else {
					record[e] = "0"
				}
			}
		}
		records = append(records, record)
	}
	return records, nil
}

// func (a *AbstractAdapter) Query(query string) (result []map[string]interface{}, err error) {

// 	rows, err := a.DB.Query(query)
// 	if err != nil {
// 		return
// 	}
// 	defer rows.Close()
// 	columns := rows.Columns()
// 	for rows.Next() {
// 		row := map[string]interface{}{}

// 		//fmt.Println(rows.Columns())
// 		err = rows.Scan(&name)
// 		if err != nil {
// 			//return
// 			log.Fatal(err)
// 		}
// 		names = append(names, name)
// 		//log.Println(name)
// 	}
// 	err = rows.Err()
// 	if err != nil {
// 		//return
// 		log.Fatal(err)
// 	}
// 	return
// }
