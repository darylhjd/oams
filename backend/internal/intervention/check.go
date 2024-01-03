package intervention

import (
	"fmt"

	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	"github.com/darylhjd/oams/backend/pkg/azmail"
	"github.com/expr-lang/expr"
	"go.uber.org/zap"
)

type checkResults struct {
	RuleFailures  map[int64][]ruleFailure                 // Map of Class IDs to rule execution failures.
	CheckFailures map[userKey][]model.ClassAttendanceRule // Map of User IDs to rule check failures.
}

type ruleFailure struct {
	Rule  model.ClassAttendanceRule
	Error error
}

func (s *Service) informRuleFailures(ruleFailures map[int64][]ruleFailure) {
	for classId, failures := range ruleFailures {
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
		RuleFailures:  map[int64][]ruleFailure{},
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
				results.RuleFailures[classId] = append(results.RuleFailures[classId], ruleFailure{
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
					results.RuleFailures[classId] = append(results.RuleFailures[classId], ruleFailure{
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
