package env

// init helps to automatically load the environment variables for a programme.
// This should be the only init function in the codebase!!!
func init() {
	MustLoad()
}
