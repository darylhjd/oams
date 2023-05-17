package env

var loaded bool

// init helps to automatically load the environment variables for a programme.
func init() {
	if !loaded {
		MustLoad()
		loaded = !loaded
	}
}
