package lib

import (
	"fmt"
	"github.com/arasuresearch/arasu/lib/stringer"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func MapToBdObj(m map[string]interface{}, i interface{}) error {
	//fmt.Println(m)
	r := reflect.ValueOf(i)
	if r.Kind() != reflect.Ptr {
		return fmt.Errorf("input argument(%s type ) must be a pointer type like (&obj)", r.Kind())
	}
	value := reflect.Indirect(r)
	for k, v := range m {
		//fmt.Println(k, v)
		//fmt.Printf("%T : %T\n", k, v)
		value.FieldByName(k).Set(reflect.ValueOf(v))
	}
	return nil
}
func MapToObj(params map[string]interface{}, obj interface{}, strict bool) error {
	r := reflect.ValueOf(obj)
	if r.Kind() != reflect.Ptr {
		return fmt.Errorf("input argument(%s type ) must be a pointer type like (&obj)", r.Kind())
	}
	elem := r.Elem()
	return mapToElem(elem, elem.Type(), params, strict)
}

func mapToElem(e reflect.Value, t reflect.Type, p map[string]interface{}, strict bool) error {
	for i := 0; i < e.NumField(); i++ {
		var sv string                 //string value
		var mv map[string]interface{} //map value
		var iv interface{}            //interface{} value
		tfield := t.Field(i)          //type Field
		key := tfield.Name

		if v, ok := p[key]; ok {
			iv = v
		} else if v, ok := p[stringer.Underscore(key)]; ok {
			iv = v
		} else {
			if tfield.Anonymous {
				afield := reflect.New(tfield.Type) //anonymous Field
				efield := afield.Elem()
				if efield.Kind() == reflect.Ptr {
					efield.Set(reflect.New(efield.Type().Elem()))
					//efield = efield.Elem()
					ele := efield.Elem()
					if e := mapToElem(ele, ele.Type(), p, strict); e != nil {
						return e
					}
					e.Field(i).Set(efield)

				} else {
					if e := mapToElem(efield, efield.Type(), p, strict); e != nil {
						return e
					}
					e.Field(i).Set(efield)
				}
				//fmt.Println(efield, efield.Type())
			}
			continue
		}
		efield := e.Field(i) //element field
		if !efield.CanSet() {
			if strict {
				return fmt.Errorf("Field (%s) is private you can't assign", key)
			}
			continue
		}

		// bool, for JSON booleans
		// float64, for JSON numbers
		// string, for JSON strings
		// []interface{}, for JSON arrays
		// map[string]interface{}, for JSON objects
		// nil for JSON null

		switch v := iv.(type) {
		case nil:
			continue
		case float64:
			sv = fmt.Sprint(v)
		case bool:
			efield.SetBool(v)
			continue
			//sv = fmt.Sprint(v)
		case string:
			sv = v
		case map[string]interface{}:
			mv = v
		default:
			return fmt.Errorf("value type is not in string or map %v", v)
		}
		typ := efield.Type().Name()
		//fmt.Println(typ)
		switch typ {
		case "int", "int8", "int16", "int32", "int64":
			if v, e := strconv.ParseInt(sv, 10, 64); e != nil {
				return e
			} else {
				efield.SetInt(v)
			}
		case "uint", "uint8", "uint16", "uint32", "uint64":
			if v, e := strconv.ParseUint(sv, 10, 64); e != nil {
				return e
			} else {
				efield.SetUint(v)
			}

		case "float32", "float64":
			if v, e := strconv.ParseFloat(sv, 64); e != nil {
				return e
			} else {
				efield.SetFloat(v)
			}
		case "bool":
			if v, e := strconv.ParseBool(sv); e != nil {
				return e
			} else {
				efield.SetBool(v)
			}
		case "[]byte":
			efield.SetBytes([]byte(sv))
		case "string":
			efield.SetString(sv)
		case "Time":
			sv = strings.Replace(sv, " ", "T", 1)
			if v, e := time.Parse(time.RFC3339Nano, sv); e != nil {
				return e
			} else {
				efield.Set(reflect.ValueOf(v))
			}
		default:
			if efield.Kind() == reflect.Ptr {
				efield.Set(reflect.New(efield.Type().Elem()))
				ele := efield.Elem()
				if e := mapToElem(ele, ele.Type(), mv, strict); e != nil {
					return e
				}
			} else {
				ele := reflect.Indirect(efield)
				if e := mapToElem(ele, ele.Type(), mv, strict); e != nil {
					return e
				}
			}
		}
	}
	return nil
}

// func mapToElem(e reflect.Value, t reflect.Type, p map[string]interface{}, strict bool) error {
// 	for i := 0; i < e.NumField(); i++ {
// 		var sv string                 //string value
// 		var mv map[string]interface{} //map value
// 		var iv interface{}            //interface{} value

// 		key := t.Field(i).Name
// 		//fmt.Println(key, lib.Camelize(key))
// 		if v, ok := p[key]; ok {
// 			iv = v
// 		} else if v, ok := p[Camelize(key)]; ok {
// 			iv = v
// 		} else {
// 			continue
// 		}
// 		if !e.Field(i).CanSet() {
// 			if strict {
// 				return fmt.Errorf("Field (%s) is private you can't assign", key)
// 			}
// 			continue
// 		}
// 		switch v := iv.(type) {
// 		// case nil:
// 		// 	continue
// 		case string:
// 			sv = v
// 		case map[string]interface{}:
// 			mv = v
// 		default:
// 			return fmt.Errorf("value type is not in string or map %v", v)
// 		}
// 		field := e.Field(i)
// 		typ := field.Type().Name()
// 		//fmt.Println(typ)
// 		switch typ {
// 		case "int", "int8", "int16", "int32", "int64":
// 			if v, e := strconv.ParseInt(sv, 10, 64); e != nil {
// 				return e
// 			} else {
// 				field.SetInt(v)
// 			}
// 		case "uint", "uint8", "uint16", "uint32", "uint64":
// 			if v, e := strconv.ParseUint(sv, 10, 64); e != nil {
// 				return e
// 			} else {
// 				field.SetUint(v)
// 			}

// 		case "float32", "float64":
// 			if v, e := strconv.ParseFloat(sv, 64); e != nil {
// 				return e
// 			} else {
// 				field.SetFloat(v)
// 			}
// 		case "bool":
// 			if v, e := strconv.ParseBool(sv); e != nil {
// 				return e
// 			} else {
// 				field.SetBool(v)
// 			}
// 		case "[]byte":
// 			field.SetBytes([]byte(sv))
// 		case "string":
// 			field.SetString(sv)
// 		case "Time":
// 			sv = strings.Replace(sv, " ", "T", 1)
// 			if v, e := time.Parse(time.RFC3339Nano, sv); e != nil {
// 				return e
// 			} else {
// 				field.Set(reflect.ValueOf(v))
// 			}
// 		default:
// 			if field.Kind() == reflect.Ptr {
// 				field.Set(reflect.New(field.Type().Elem()))
// 				ele := field.Elem()
// 				if e := mapToElem(ele, ele.Type(), mv, strict); e != nil {
// 					return e
// 				}
// 			} else {
// 				ele := reflect.Indirect(field)
// 				if e := mapToElem(ele, ele.Type(), mv, strict); e != nil {
// 					return e
// 				}
// 			}
// 		}
// 	}
// 	return nil
// }

// func getMapValue(params map[string]interface{}, key string) interface{} {
// 	//fmt.Println(params, key)
// 	return params[key]
// }

// func SetStruct(obj interface{}, params map[string]interface{}) error {
// 	r := reflect.ValueOf(obj)
// 	if r.Kind() != reflect.Ptr {
// 		return fmt.Errorf("input argument(%s type ) must be a pointer type like (&obj)", r.Kind())
// 	}
// 	e := r.Elem()
// 	t := e.Type()
// 	t0 := e.Field(0)

// 	if t0.Kind() == reflect.Ptr {
// 		r0 := reflect.New(t0.Type().Elem())
// 		t0.Set(r0)
// 		e0 := t0.Elem()
// 		e0.Field(0).SetString("lollusabha")
// 	} else {
// 		r2 := reflect.Indirect(t0)
// 		r2.Field(0).SetString("lollusabha")

// 	}
// 	for i := 0; i < e.NumField(); i++ {
// 		var sv string                 //string value
// 		var mv map[string]interface{} //map value
// 		//fmt.Println(t.Field(i).Name)
// 		if v := getMapValue(params, t.Field(i).Name); v == nil {
// 			continue
// 		} else if val, ok := v.(string); ok {
// 			sv = val
// 		} else if val, ok := v.(map[string]interface{}); ok {
// 			mv = val
// 		} else {
// 			return fmt.Errorf("value type is not in sting or map %v", v)
// 		}
// 		_ = mv
// 		field := e.Field(i)
// 		typ := field.Type().Name()
// 		//fmt.Println(typ)

// 		switch typ {
// 		case "int", "int8", "int16", "int32", "int64":
// 			if v, e := strconv.ParseInt(sv, 10, 64); e != nil {
// 				return e
// 			} else {
// 				field.SetInt(v)
// 			}
// 		case "uint", "uint8", "uint16", "uint32", "uint64":
// 			if v, e := strconv.ParseUint(sv, 10, 64); e != nil {
// 				return e
// 			} else {
// 				field.SetUint(v)
// 			}

// 		case "float32", "float64":
// 			if v, e := strconv.ParseFloat(sv, 64); e != nil {
// 				return e
// 			} else {
// 				field.SetFloat(v)
// 			}
// 		case "bool":
// 			if v, e := strconv.ParseBool(sv); e != nil {
// 				return e
// 			} else {
// 				field.SetBool(v)
// 			}
// 		case "[]byte":
// 			field.SetBytes([]byte(sv))
// 		case "string":
// 			field.SetString(sv)
// 		case "Time":
// 			sv = strings.Replace(sv, " ", "T", 1)
// 			if v, e := time.Parse(time.RFC3339Nano, sv); e != nil {
// 				return e
// 			} else {
// 				field.Set(reflect.ValueOf(v))
// 			}
// 		default:
// 			fmt.Println(typ)
// 		}
// 	}

// 	return nil
// }
