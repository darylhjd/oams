package intervention

import (
	"fmt"

	"github.com/expr-lang/expr"
	"go.uber.org/zap"
)

func (s *Service) performChecks(fGroup factGrouping, rGroup ruleGrouping) error {
	// For each class rule, check which user fails the rule.
	for classId, students := range fGroup {
		s.l.Info(fmt.Sprintf("%s - performing rule checks", Namespace), zap.Int64("class_id", classId))

		for _, rule := range rGroup[classId] {
			s.l.Info(
				fmt.Sprintf("%s - checking rule", Namespace),
				zap.String("title", rule.Title),
				zap.String("description", rule.Description),
			)

			prg, err := expr.Compile(rule.Rule, expr.AsBool(), expr.Env(rule.Environment.Env))
			if err != nil {
				return err
			}

			for studentId, studentFacts := range students {
				runEnv := rule.Environment.Env.SetFacts(studentFacts)

				res, err := expr.Run(prg, runEnv)
				if err != nil {
					return err
				}

				if res.(bool) {
					s.l.Info(
						fmt.Sprintf("%s - student failed check", Namespace),
						zap.String("user_id", studentId),
					)
				}
			}
		}
	}

	return nil
}
