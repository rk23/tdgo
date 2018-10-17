package fields

import "sort"

func Append(fields ...string) string {
	sort.Strings(fields)
	a := ""
	for _, field := range fields {
		a += field + ","
	}
	return a
}
