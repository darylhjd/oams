package intervention

import (
	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	"github.com/darylhjd/oams/backend/internal/intervention/fact"
)

type factGrouping map[int64]map[string][]fact.F

// groupFacts first by class, and then by user.
func (s *Service) groupFacts(facts []fact.F) factGrouping {
	grouping := factGrouping{}

	for _, f := range facts {
		if _, ok := grouping[f.ClassID]; !ok {
			grouping[f.ClassID] = map[string][]fact.F{}
		} else if _, ok = grouping[f.ClassID][f.UserID]; !ok {
			grouping[f.ClassID][f.UserID] = []fact.F{}
		}

		grouping[f.ClassID][f.UserID] = append(grouping[f.ClassID][f.UserID], f)
	}

	return grouping
}

type ruleGrouping map[int64][]model.ClassAttendanceRule

// groupRules by class.
func (s *Service) groupRules(rules []model.ClassAttendanceRule) ruleGrouping {
	grouping := ruleGrouping{}

	for _, rule := range rules {
		grouping[rule.ClassID] = append(grouping[rule.ClassID], rule)
	}

	return grouping
}
