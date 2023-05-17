package env

import "log"

var loaded bool

// init helps to automatically load the environment variables for a programme.
func init() {
	if !loaded {
		log.Print("Loading environment variables...")
		MustLoad()
		log.Println("Loaded!")
		loaded = !loaded
	}
}
