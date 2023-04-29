/*
Package autoloader is used to load the environment variables automatically by importing this
library for side effects. Only in certain circumstances should this be used:
- When creating a new cmd.
- At least once within the `tests` package (which is already done!).
Please use this wisely! It is tempting to import this everywhere but this will result in
unnecessary calls to this function.
*/
package autoloader

import (
	"github.com/darylhjd/oats/backend/env"
)

var loaded bool

func init() {
	if !loaded {
		env.MustLoad()
		loaded = !loaded
	}
}
