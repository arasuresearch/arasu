package lib

import (
	"strings"
)

func ParseDoubleTemplateString(s string, data interface{}) (string, error) {
	localStr := forwardEvenReplace(s, "{{", "_ocb_")
	localStr = backwardOddReplace(localStr, "}}", "_ccb_")
	ParsedStr, err := ParseAndExecuteTemplateText(localStr, data)
	if err != nil {
		return ParsedStr, err
	}
	ParsedStr = strings.Replace(ParsedStr, "_ocb_", "{{", -1)
	ParsedStr = strings.Replace(ParsedStr, "_ccb_", "}}", -1)
	return ParsedStr, nil
}
func forwardEvenReplace(s, sep, rep string) (r string) {
	a := strings.Split(s, sep)
	for i := 0; i < len(a); i++ {
		if i%2 == 0 {
			r = r + a[i] + rep
		} else {
			r = r + a[i] + sep
		}
	}
	r = strings.TrimSuffix(r, rep)
	return
}
func backwardOddReplace(s, sep, rep string) (r string) {
	a := strings.Split(s, sep)
	for i := 0; i < len(a); i++ {
		if i%2 == 0 {
			r = r + a[i] + sep
		} else {
			r = r + a[i] + rep
		}
	}
	r = strings.TrimSuffix(r, sep)
	return
}
