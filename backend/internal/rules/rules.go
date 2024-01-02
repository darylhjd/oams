package rules

import (
	_ "embed"
	"errors"

	"github.com/expr-lang/expr"
)

var (
	//go:embed missed_consecutive_classes.expr
	consecutiveRule string

	//go:embed min_percentage_attendance_from_session.expr
	percentageRule string
)

type RuleParams struct {
	Title             string            `json:"title"`
	Description       string            `json:"description"`
	RuleType          T                 `json:"rule_type"`
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

func (r RuleParams) Verify() (rule string, env E, err error) {
	if len(r.Title) == 0 {
		return "", nil, errors.New("title is empty")
	}

	if len(r.Description) == 0 {
		return "", nil, errors.New("description is empty")
	}

	switch r.RuleType {
	case TConsecutive:
		return r.verifyConsecutiveRule()
	case TPercentage:
		return r.verifyPercentageRule()
	case TAdvanced:
		return r.verifyAdvancedRule()
	default:
		return "", nil, errors.New("unknown rule type")
	}
}

func (r RuleParams) verifyConsecutiveRule() (rule string, env E, err error) {
	if r.ConsecutiveParams.ConsecutiveClasses < 1 {
		return "", nil, errors.New("number of consecutive classes cannot be less than 1")
	}

	env = ConsecutiveE{
		BaseE: BaseE{
			EnvType: TConsecutive,
		},
		ConsecutiveClasses: r.ConsecutiveParams.ConsecutiveClasses,
	}

	_, err = expr.Compile(consecutiveRule, expr.AsBool(), expr.Env(env))
	return consecutiveRule, env, err
}

func (r RuleParams) verifyPercentageRule() (rule string, env E, err error) {
	params := r.PercentageParams

	if params.Percentage < 0 {
		return "", nil, errors.New("percentage cannot be negative")
	} else if params.Percentage > 100 {
		return "", nil, errors.New("percentage cannot be more than 100")
	}

	if params.FromSession < 1 {
		return "", nil, errors.New("number of sessions cannot be less than 1")
	}

	env = PercentageE{
		BaseE: BaseE{
			EnvType: TPercentage,
		},
		Percentage:  params.Percentage,
		FromSession: params.FromSession,
	}

	_, err = expr.Compile(percentageRule, expr.AsBool(), expr.Env(env))
	return percentageRule, env, err
}

func (r RuleParams) verifyAdvancedRule() (rule string, env E, err error) {
	env = BaseE{
		EnvType: TAdvanced,
	}

	_, err = expr.Compile(r.AdvancedParams.Rule, expr.AsBool(), expr.Env(env))
	return r.AdvancedParams.Rule, env, err
}
