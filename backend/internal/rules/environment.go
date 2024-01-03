package rules

// T is an identifier for the type of environment.
type T int

const (
	TConsecutive T = iota
	TPercentage
	TAdvanced
)

// E is an interface that all environments for any rule must satisfy.
type E interface {
	SetFacts([]Fact) E
}

// BaseE represents the base environment variables that all environments must have.
type BaseE struct {
	EnvType     T      `json:"env_type"`
	Enrollments []Fact `expr:"enrollments" json:"-"`
}

func (e BaseE) SetFacts(facts []Fact) E {
	e.Enrollments = facts
	return e
}

type ConsecutiveE struct {
	BaseE
	ConsecutiveClasses int `expr:"consecutive_classes" json:"consecutive_classes"`
}

func (e ConsecutiveE) SetFacts(facts []Fact) E {
	e.Enrollments = facts
	return e
}

type PercentageE struct {
	BaseE
	Percentage  float64 `expr:"percentage" json:"percentage"`
	FromSession int     `expr:"from_session" json:"from_session"`
}

func (e PercentageE) SetFacts(facts []Fact) E {
	e.Enrollments = facts
	return e
}
