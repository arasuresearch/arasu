package lib

import (
	"strings"
	"time"
)

// func Constantize(s string) string {
// 	s = strings.Join(strings.Split(s, "_"), " ")
// 	s = strings.Title(s)
// 	s = strings.Join(strings.Split(s, " "), "")
// 	return s
// }
// func Camelize(s string) string {
// 	b := make([]byte, 0, len(s)*2)
// 	for i, e := range s {
// 		if e > 64 && e < 91 {
// 			if i == 0 {
// 				b = append(b, byte(e+32))
// 			} else {
// 				b = append(b, byte(rune('_')), byte(e+32))
// 			}
// 		} else {
// 			b = append(b, byte(e))
// 		}
// 	}
// 	return string(b)
// }
// func Titleize(s string) string {
// 	var r string
// 	for _, e := range strings.Split(s, "_") {
// 		r = r + strings.Title(e)
// 	}
// 	return r
// }

// // below constantice has error
// // Contact #ontact

// func Constantize(s string) string {
// 	b := make([]byte, 0, len(s))
// 	var flag = true
// 	for _, e := range s {
// 		if e == 95 {
// 			flag = true
// 			continue
// 		} else {
// 			if flag {
// 				flag = false
// 				b = append(b, byte(e-32))
// 			} else {
// 				b = append(b, byte(e))
// 			}
// 		}
// 	}
// 	return string(b)
// }

var HiddenFuncs = []string{"Open", "Serve", "Close", "RedirectTo", "Render", "BeforeFunc", "AfterFunc"}

func IsItVisible(name string) bool {
	for i := 0; i < len(HiddenFuncs); i++ {
		if HiddenFuncs[i] == name {
			return false
		}
	}
	return true

}

func (s *AS) Exists(str string) bool {
	for _, e := range *s {
		if e == str {
			return true
		}
	}
	return false
}
func (s *AS) Add(str string) bool {
	if !s.Exists(str) {
		*s = append(*s, str)
		return true
	}
	return false
}
func (s *AS) Del(str string) (r bool) {
	for i, e := range *s {
		if e == str {
			(*s)[i] = (*s)[len(*s)-1]
			(*s) = (*s)[0 : len(*s)-1]
			r = true
		}
	}
	return r
}

func ParseKeyValueAndRemaningArguments(cmd_args []string) (map[string]string, []string) {
	kv := map[string]string{}
	args := []string{}

	for _, e := range cmd_args {
		if strings.Contains(e, ":") {
			key_val := strings.SplitN(e, ":", 2)
			kv[key_val[0]] = key_val[1]
		} else {
			args = append(args, e)

		}
	}
	return kv, args
}

func UniqueTimeVersion() string {
	s := strings.Split(time.Now().String(), ".")[0]
	s = strings.Replace(s, ":", "", -1)
	s = strings.Replace(s, "-", "", -1)
	s = strings.Replace(s, " ", "", -1)
	return s
}
