package rules

import (
	_ "embed"
	"errors"

	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	"github.com/expr-lang/expr"
)

var (
	//go:embed missed_consecutive_classes.expr
	missedConsecutiveClasses string

	//go:embed min_percentage_attendance_from_session.expr
	minPercentageAttendanceFromSession string
)

type RuleType string

const (
	RuleTypeSimple   RuleType = "simple"
	RuleTypeAdvanced RuleType = "advanced"
)

type SimpleRule string

const (
	SimpleRuleMissedConsecutiveClasses           SimpleRule = "missed_consecutive_classes"
	SimpleRuleMinPercentageAttendanceFromSession SimpleRule = "min_percentage_attendance_from_session"
)

type baseEnv struct {
	Enrollments []model.SessionEnrollment `expr:"enrollments"`
}

type missedConsecutiveClassesEnv struct {
	baseEnv
	ConsecutiveClasses int `expr:"consecutive_classes"`
}

type minPercentageAttendanceFromSessionEnv struct {
	baseEnv
	Percentage float64 `expr:"percentage"`
	From       int     `expr:"from"`
}

var baseEnvironment = baseEnv{
	Enrollments: []model.SessionEnrollment{},
}

type RuleParams struct {
	Title             string     `json:"title"`
	Description       string     `json:"description"`
	RuleType          RuleType   `json:"rule_type"`
	PresetRule        SimpleRule `json:"preset_rule"`
	ConsecutiveParams struct {
		Num int `json:"num"`
	} `json:"consecutive_params"`
	PercentageParams struct {
		Percentage float64 `json:"percentage"`
		From       int     `json:"from"`
	} `json:"percentage_params"`
	Rule string `json:"rule"`
}

func (r RuleParams) Verify() error {
	if len(r.Title) == 0 {
		return errors.New("title is empty")
	}

	if len(r.Description) == 0 {
		return errors.New("description is empty")
	}

	switch r.RuleType {
	case RuleTypeSimple:
		return r.verifySimpleRule()
	case RuleTypeAdvanced:
		return r.verifyAdvancedRule()
	default:
		return errors.New("unknown rule type")
	}
}

func (r RuleParams) verifySimpleRule() error {
	switch r.PresetRule {
	case SimpleRuleMissedConsecutiveClasses:
		_, err := expr.Compile(missedConsecutiveClasses, expr.AsBool(), expr.Env(missedConsecutiveClassesEnv{
			baseEnvironment,
			r.ConsecutiveParams.Num,
		}))
		return err
	case SimpleRuleMinPercentageAttendanceFromSession:
		_, err := expr.Compile(minPercentageAttendanceFromSession, expr.AsBool(), expr.Env(minPercentageAttendanceFromSessionEnv{
			baseEnvironment,
			r.PercentageParams.Percentage,
			r.PercentageParams.From,
		}))
		return err
	default:
		return errors.New("unknown simple rule")
	}
}

func (r RuleParams) verifyAdvancedRule() error {
	_, err := expr.Compile(r.Rule, expr.AsBool(), expr.Env(baseEnvironment))
	return err
}
