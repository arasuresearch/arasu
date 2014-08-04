package lib

type H map[interface{}]interface{}
type HS map[string]interface{}
type HSS map[string]string

type A []interface{}
type AS []string
type AI []int

type I interface{}

func (h H) Get(key string) I {
	if v, ok := h[key]; ok {
		return v
	} else {
		return nil
	}
}

// func SetStruct(s interface{}, args map[string]interface{}) {
// 	r := reflect.ValueOf(s).Elem()
// 	t := r.Type()
// 	for k, v := range args {
// 		if _, ok := t.FieldByName(k); ok {
// 			f0 := r.FieldByName(k)
// 			vv := reflect.ValueOf(v)
// 			f0.Set(vv)
// 		}
// 	}
// }

// func MapToStruct(m map[string]interface{}, s interface{}) error {
// 	tmp, err := json.Marshal(m)
// 	if err != nil {
// 		return err
// 	}
// 	err = json.Unmarshal(tmp, s)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
