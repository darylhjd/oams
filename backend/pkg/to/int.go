package to

import "strconv"

// Int64 converts a string number in base 10 to an int64 type.
func Int64(n string) (int64, error) {
	return strconv.ParseInt(n, 10, 64)
}
