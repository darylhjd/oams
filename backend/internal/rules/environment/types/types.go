package types

// T is an identifier for the type of environment.
type T int

const (
	TConsecutive T = iota
	TPercentage
	TAdvanced
)
