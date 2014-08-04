package lib

import (
	"strings"
)

func ArgsWithoutOptions(a []string) []string {
	for i, e := range a {
		if strings.HasPrefix(e, "-") {
			return a[:i]
		}
	}
	return a
}
