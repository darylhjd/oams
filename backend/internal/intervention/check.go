package intervention

import (
	"fmt"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/expr-lang/expr"
	"go.uber.org/zap"
)

// userFailedRules contains lists of database.RuleInfo for each user.
// The list consists of database.RuleInfo that each user broke.
type userFailedRules map[userKey][]database.RuleInfo

// ruleCreatorRuleFailedUsers contains lists of ruleAndFailedUsers for each rule creator.
// The list consists of ruleAndFailedUsers, of which the rule belongs to the creator.
type ruleCreatorRuleFailedUsers map[userKey][]ruleAndFailedUsers

func (s *Service) performChecks(fGroup factGrouping, rGroup ruleGrouping) (userFailedRules, ruleCreatorRuleFailedUsers, error) {
	userRules := userFailedRules{}
	ruleCreatorUsers := ruleCreatorRuleFailedUsers{}

	for classId, group := range fGroup {
		for _, rule := range rGroup[classId] {
			s.l.Info(
				fmt.Sprintf("%s - performing rule checks", Namespace),
				zap.Int64("class_id", classId),
				zap.Int64("rule_id", rule.ID),
			)

			prg, err := expr.Compile(rule.Rule, expr.AsBool(), expr.Env(rule.Environment.Env))
			if err != nil {
				return nil, nil, fmt.Errorf("failed to compile rule with id %d: %w", rule.ID, err)
			}

			ruleFailedUsers := ruleAndFailedUsers{Rule: rule}
			for user, facts := range group {
				runEnv := rule.Environment.Env.SetFacts(facts)

				res, err := expr.Run(prg, runEnv)
				switch {
				case err != nil:
					return nil, nil, fmt.Errorf("failed to run rule with id %d: %w", rule.ID, err)
				case res.(bool):
					userRules[user] = append(userRules[user], rule)
					ruleFailedUsers.FailedUsers = append(ruleFailedUsers.FailedUsers, user)
				}
			}

			creatorKey := userKey{rule.CreatorID, rule.CreatorName, rule.CreatorEmail}
			ruleCreatorUsers[creatorKey] = append(ruleCreatorUsers[creatorKey], ruleFailedUsers)
		}
	}

	return userRules, ruleCreatorUsers, nil
}
