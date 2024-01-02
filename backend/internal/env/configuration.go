package env

import (
	"os"
)

const (
	configuration = "CONFIGURATION"
)

type Conf string

const (
	ConfAPIServer    Conf = "apiserver"
	ConfIntervention Conf = "intervention"
)

// GetConfiguration returns the CONFIGURATION environment variable.
func GetConfiguration() Conf {
	return Conf(os.Getenv(configuration))
}
