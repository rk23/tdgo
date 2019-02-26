package fields

import (
	"sort"
	"strconv"
)

// Append takes in a variable number of string arguments and adds them
// to a single string in order
func Append(fields ...int) string {
	sort.Ints(fields)
	a := ""
	for idx, field := range fields {
		a += strconv.Itoa(field)
		if idx < len(fields)-1 {
			a += ","
		}
	}
	return a
}
