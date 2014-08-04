package lib

import (
	"strings"
)

func ParseFlag(names []string, value string, args []string) string {
	if len(names) == 0 {
		return value
	}
	for _, e := range names {
		fopt := AS{"--" + e, "-" + e}
		for i, e := range args {
			if fopt.Exists(e) {
				if len(args) > i+1 {
					return args[i+1]
				}
			}
			if strings.HasPrefix(e, fopt[0]+"=") {
				t := strings.TrimPrefix(e, fopt[0]+"=")
				if len(t) > 0 {
					return t
				}
			}
			if strings.HasPrefix(e, fopt[1]+"=") {
				t := strings.TrimPrefix(e, fopt[1]+"=")
				if len(t) > 0 {
					return t
				}
			}

		}
	}
	return value
}
