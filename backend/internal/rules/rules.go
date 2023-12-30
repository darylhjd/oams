package rules

import (
	_ "embed"
	"errors"

	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	"github.com/expr-lang/expr"
)

var (
	//go:embed missed_consecutive_classes.expr
	consecutiveRule string

	//go:embed min_percentage_attendance_from_session.expr
	percentageRule string
)

type RuleType string

const (
	RuleTypeMissedConsecutiveClasses           RuleType = "missed_consecutive_classes"
	RuleTypeMinPercentageAttendanceFromSession RuleType = "min_percentage_attendance_from_session"
	RuleTypeAdvanced                           RuleType = "advanced"
)

type baseEnvironment struct {
	Enrollments []model.SessionEnrollment `expr:"enrollments"`
}

var baseEnv = baseEnvironment{
	Enrollments: []model.SessionEnrollment{},
}

type missedConsecutiveClassesEnv struct {
	baseEnvironment
	ConsecutiveClasses int `expr:"consecutive_classes"`
}

type minPercentageAttendanceFromSessionEnv struct {
	baseEnvironment
	Percentage  float64 `expr:"percentage"`
	FromSession int     `expr:"from_session"`
}

type RuleParams struct {
	Title             string            `json:"title"`
	Description       string            `json:"description"`
	RuleType          RuleType          `json:"rule_type"`
	ConsecutiveParams consecutiveParams `json:"consecutive_params"`
	PercentageParams  percentageParams  `json:"percentage_params"`
	AdvancedParams    advancedParams    `json:"advanced_params"`
}

type consecutiveParams struct {
	ConsecutiveClasses int `json:"consecutive_classes"`
}

type percentageParams struct {
	Percentage  float64 `json:"percentage"`
	FromSession int     `json:"from_session"`
}

type advancedParams struct {
	Rule string `json:"rule"`
}

func (r RuleParams) Verify() (rule string, env any, err error) {
	if len(r.Title) == 0 {
		return "", nil, errors.New("title is empty")
	}

	if len(r.Description) == 0 {
		return "", nil, errors.New("description is empty")
	}

	switch r.RuleType {
	case RuleTypeMissedConsecutiveClasses:
		return r.verifyConsecutiveRule()
	case RuleTypeMinPercentageAttendanceFromSession:
		return r.verifyPercentageRule()
	case RuleTypeAdvanced:
		return r.verifyAdvancedRule()
	default:
		return "", nil, errors.New("unknown rule type")
	}
}

func (r RuleParams) verifyConsecutiveRule() (rule string, env any, err error) {
	if r.ConsecutiveParams.ConsecutiveClasses < 1 {
		return "", nil, errors.New("number of consecutive classes cannot be less than 1")
	}

	env = missedConsecutiveClassesEnv{
		baseEnv,
		r.ConsecutiveParams.ConsecutiveClasses,
	}

	_, err = expr.Compile(consecutiveRule, expr.AsBool(), expr.Env(env))
	return consecutiveRule, env, err
}

func (r RuleParams) verifyPercentageRule() (rule string, env any, err error) {
	params := r.PercentageParams

	if params.Percentage < 0 {
		return "", nil, errors.New("percentage cannot be negative")
	}

	if params.FromSession < 1 {
		return "", nil, errors.New("number of sessions cannot be less than 1")
	}

	env = minPercentageAttendanceFromSessionEnv{
		baseEnv,
		params.Percentage,
		params.FromSession,
	}

	_, err = expr.Compile(percentageRule, expr.AsBool(), expr.Env(env))
	return percentageRule, env, err
}

func (r RuleParams) verifyAdvancedRule() (rule string, env any, err error) {
	_, err = expr.Compile(r.AdvancedParams.Rule, expr.AsBool(), expr.Env(baseEnv))
	return r.AdvancedParams.Rule, baseEnv, err
}
