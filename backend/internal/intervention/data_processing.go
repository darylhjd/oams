package intervention

import (
	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/intervention/fact"
)

type factGrouping map[int64]map[userKey][]fact.F

type userKey struct {
	ID    string
	Name  string
	Email string
}

// groupFacts first by class, and then by user.
func (s *Service) groupFacts(facts []fact.F) factGrouping {
	grouping := factGrouping{}

	for _, f := range facts {
		if _, ok := grouping[f.ClassID]; !ok {
			grouping[f.ClassID] = map[userKey][]fact.F{}
		}

		key := userKey{
			ID:    f.UserID,
			Name:  f.UserName,
			Email: f.UserEmail,
		}

		if _, ok := grouping[f.ClassID][key]; !ok {
			grouping[f.ClassID][key] = []fact.F{}
		}

		grouping[f.ClassID][key] = append(grouping[f.ClassID][key], f)
	}

	return grouping
}

type ruleGrouping map[int64][]database.RuleInfo

// groupRules by class.
func (s *Service) groupRules(rules []database.RuleInfo) ruleGrouping {
	grouping := ruleGrouping{}

	for _, rule := range rules {
		grouping[rule.ClassID] = append(grouping[rule.ClassID], rule)
	}

	return grouping
}
