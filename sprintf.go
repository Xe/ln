package ln

import "fmt"

// Fmt formats a string like sprintf.
func Fmt(format string, args ...interface{}) F {
	return F{
		"msg": fmt.Sprintf(format, args...),
	}
}
