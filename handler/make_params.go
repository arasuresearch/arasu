package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

//buep base_url_embeded_params
func parseParams(r *http.Request, buep map[string]string) (params map[string]interface{}, err error) {

	if err = r.ParseForm(); err != nil {
		return
	}

	if strings.Contains(r.UserAgent(), "(Dart)") {
		params, err = jsonParams(r)
	} else {
		ctype := r.Header.Get("Content-Type")
		//fmt.Println(ctype)
		switch {
		case strings.Contains(ctype, "application/json"):
			data := make([]byte, r.ContentLength)
			if _, err = r.Body.Read(data); err != nil {
				return
			}
			var values url.Values
			values, err = url.ParseQuery(string(data))
			if err != nil {
				return
			}
			params, err = encodedParams(values)
		case strings.Contains(ctype, "application/x-www-form-urlencoded"):
			params, err = encodedParams(r.Form)
		// case strings.Contains(ctype, "multipart/form-data"):
		// 	params, err = encodedParams(r.Form)
		// case strings.Contains(ctype, "application/xml"):
		// 	params, err = encodedParams(r.Form)
		default:
			params, err = encodedParams(r.Form)
		}
	}
	if err != nil {
		return
	}
	for k, v := range buep {
		params[k] = v
	}
	return
}

func jsonParams(r *http.Request) (map[string]interface{}, error) {
	if r.ContentLength < 1 {
		return map[string]interface{}{}, nil
	}
	data := make([]byte, r.ContentLength)

	if _, err := r.Body.Read(data); err != nil && err != io.EOF {
		return nil, err
	}

	//fmt.Println(string(data))

	var params map[string]interface{}
	if err := json.Unmarshal(data, &params); err != nil {
		return nil, err
	}
	return params, nil
}
func encodedParams(Values map[string][]string) (map[string]interface{}, error) {
	r := map[string]interface{}{}
	for k, v := range Values {
		var value interface{}
		if len(v) == 1 {
			value = v[0]
		} else {
			return nil, fmt.Errorf("%s has multiple[%d] values (%v)", k, len(v), v)
		}
		if strings.Contains(k, "]") {
			strs := strings.SplitN(strings.Replace(k, "]", "", -1), "[", -1)
			lastKey := strs[len(strs)-1]
			rr := r
			strs = strs[:len(strs)-1]
			for _, e := range strs {
				if rr[e] == nil {
					rr[e] = make(map[string]interface{})
				}
				rr = rr[e].(map[string]interface{})
			}
			rr[lastKey] = value
		} else {
			r[k] = value
		}
	}
	return r, nil
}

func parseArgsFromParams(t reflect.Type, args []string, p map[string]interface{}) ([]reflect.Value, error) {
	if len(args) == 0 {
		return nil, nil
	}
	res := make([]reflect.Value, t.NumIn())
	for i := 0; i < t.NumIn(); i++ {
		kind := t.In(i).Kind().String()
		inp := t.In(i).String()
		switch kind {
		case "map":
			if reflectValue, err := makeMapParam(args[i], inp, p); err != nil {
				return res, err
			} else {
				res[i] = reflectValue
			}
		case "array", "slice":
			if reflectValue, err := makeArrayParam(args[i], inp, p); err != nil {
				return res, err
			} else {
				res[i] = reflectValue
			}

		default:
			if reflectValue, err := makeBuiltinParam(args[i], inp, p); err != nil {
				return res, err
			} else {
				res[i] = reflectValue
			}
		}
	}
	return res, nil
}
func makeArrayParam(name, typ string, p map[string]interface{}) (rv reflect.Value, err error) {
	var value interface{}
	if v, ok := p[name]; ok {
		value = v
	} else {
		return rv, fmt.Errorf("argument '%s' is not set by client ", name)
	}
	switch typ {
	case "[]string":
		if res, ok := value.([]string); ok {
			return reflect.ValueOf(res), nil
		} else {
			return rv, fmt.Errorf("argument '%s' is not matching the expected type %s ", name, typ)
		}
	}
	return rv, fmt.Errorf("argument '%s' format '%s' is not matching the basic types (like int,string,map slice,etc)", name, typ)
}

func makeMapParam(name, typ string, p map[string]interface{}) (rv reflect.Value, err error) {
	var value map[string]interface{}
	if v, ok := p[name]; ok {
		if v0, ok0 := v.(map[string]interface{}); ok0 {
			value = v0
		} else {
			return rv, fmt.Errorf("argument '%s' is not matching the expected type %s ", name, typ)
		}
	} else {
		return rv, fmt.Errorf("argument '%s' is not set by client ", name)
	}

	switch typ {
	case "map[string]string":
		res := map[string]string{}
		for k, v := range value {
			res[k] = v.(string)
		}
		return reflect.ValueOf(res), nil
	case "map[string]interface {}":
		return reflect.ValueOf(value), nil
	}
	return rv, fmt.Errorf("argument '%s' format '%s' is not matching the basic types (like int,string,map slice,etc)", name, typ)
}

func makeBuiltinParam(name, typ string, p map[string]interface{}) (rv reflect.Value, err error) {
	var value string
	if v, ok := p[name]; ok {
		if v0, ok0 := v.(string); ok0 {
			value = v0
		} else {
			return rv, fmt.Errorf("argument '%s' is not matching the expected type %s ", name, typ)
		}
	} else {
		return rv, fmt.Errorf("argument '%s' is not set by client ", name)
	}

	switch typ {
	case "int":
		if v, err := strconv.ParseInt(value, 10, 0); err != nil {
			return rv, err
		} else {
			return reflect.ValueOf(int(v)), nil
		}
	case "int8":
		if v, err := strconv.ParseInt(value, 10, 8); err != nil {
			return rv, err
		} else {
			return reflect.ValueOf(int8(v)), nil
		}

	case "int16":
		if v, err := strconv.ParseInt(value, 10, 16); err != nil {
			return rv, err
		} else {
			return reflect.ValueOf(int16(v)), nil
		}
	case "int32":
		if v, err := strconv.ParseInt(value, 10, 32); err != nil {
			return rv, err
		} else {
			return reflect.ValueOf(int32(v)), nil
		}
	case "int64":
		if v, err := strconv.ParseInt(value, 10, 64); err != nil {
			return rv, err
		} else {
			return reflect.ValueOf(int64(v)), nil
		}
	case "uint":
		if v, err := strconv.ParseUint(value, 10, 0); err != nil {
			return rv, err
		} else {
			return reflect.ValueOf(uint(v)), nil
		}
	case "uint8":
		if v, err := strconv.ParseUint(value, 10, 8); err != nil {
			return rv, err
		} else {
			return reflect.ValueOf(uint8(v)), nil
		}

	case "uint16":
		if v, err := strconv.ParseUint(value, 10, 16); err != nil {
			return rv, err
		} else {
			return reflect.ValueOf(uint16(v)), nil
		}
	case "uint32":
		if v, err := strconv.ParseUint(value, 10, 32); err != nil {
			return rv, err
		} else {
			return reflect.ValueOf(uint32(v)), nil
		}
	case "uint64":
		if v, err := strconv.ParseUint(value, 10, 64); err != nil {
			return rv, err
		} else {
			return reflect.ValueOf(uint64(v)), nil
		}
	case "float32":
		if v, err := strconv.ParseFloat(value, 32); err != nil {
			return rv, err
		} else {
			return reflect.ValueOf(float32(v)), nil
		}
	case "float64":
		if v, err := strconv.ParseFloat(value, 64); err != nil {
			return rv, err
		} else {
			return reflect.ValueOf(float64(v)), nil
		}
	case "bool":
		if v, err := strconv.ParseBool(value); err != nil {
			return rv, err
		} else {
			return reflect.ValueOf(v), nil
		}
	case "string":
		return reflect.ValueOf(value), nil
	default:
		//return reflect.ValueOf(value), nil
	}
	return rv, fmt.Errorf("format '%s' is not matching the basic types (like int,string,map slice,etc)", typ)
}
