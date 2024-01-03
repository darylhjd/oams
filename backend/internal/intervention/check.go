package intervention

import (
	"fmt"

	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	"github.com/darylhjd/oams/backend/pkg/azmail"
	"github.com/expr-lang/expr"
	"go.uber.org/zap"
)

type checkResults struct {
	RuleFailures  ruleFailures
	CheckFailures checkFailures
}

// ruleFailures is a map of Class IDs to rule execution failures.
type ruleFailures map[int64][]ruleError

type ruleError struct {
	Rule  model.ClassAttendanceRule
	Error error
}

// checkFailures is a map of User IDs to failed rules.
type checkFailures map[userKey][]model.ClassAttendanceRule

func (s *Service) informRuleFailures(rFailures ruleFailures) {
	for classId, failures := range rFailures {
		for _, failure := range failures {
			s.l.Error(fmt.Sprintf("%s - rule failed to execute", Namespace),
				zap.Int64("class_id", classId),
				zap.Int64("rule_id", failure.Rule.ID),
				zap.Error(failure.Error),
			)
		}
	}
}

func (s *Service) performChecks(fGroup factGrouping, rGroup ruleGrouping) checkResults {
	results := checkResults{
		RuleFailures:  ruleFailures{},
		CheckFailures: map[userKey][]model.ClassAttendanceRule{},
	}

	for classId, users := range fGroup {
		s.l.Info(fmt.Sprintf("%s - performing rule checks", Namespace), zap.Int64("class_id", classId))

		for _, rule := range rGroup[classId] {
			prg, err := expr.Compile(rule.Rule, expr.AsBool(), expr.Env(rule.Environment.Env))
			if err != nil {
				s.l.Error(
					fmt.Sprintf("%s - error compiling rule", Namespace),
					zap.Int64("rule_id", rule.ID),
				)
				results.RuleFailures[classId] = append(results.RuleFailures[classId], ruleError{
					rule, err,
				})
				continue
			}

			for user, userFacts := range users {
				runEnv := rule.Environment.Env.SetFacts(userFacts)

				res, err := expr.Run(prg, runEnv)
				if err != nil {
					s.l.Error(
						fmt.Sprintf("%s - error running rule for student", Namespace),
						zap.Int64("rule_id", rule.ID),
						zap.String("user_id", user.ID),
					)
					results.RuleFailures[classId] = append(results.RuleFailures[classId], ruleError{
						rule, err,
					})
					continue
				}

				if res.(bool) {
					results.CheckFailures[user] = append(results.CheckFailures[user], rule)
				}
			}
		}
	}

	return results
}

func (s *Service) processCheckResults(results checkResults) ([]*azmail.Mail, error) {
	s.informRuleFailures(results.RuleFailures)
	return s.generateMails(results.CheckFailures)
}
