package intervention

import (
	"strings"

	"github.com/darylhjd/oams/backend/pkg/azmail"
)

// ruleCreatorEmailGrouping is a map of rule creator ID to their email arguments.
// This is a helper type to help generate the arguments required.
type ruleCreatorEmailGrouping map[string]ruleCreatorEmailArgs

func (s *Service) generateMails(cFailures checkFailures) ([]*azmail.Mail, error) {
	mails := make([]*azmail.Mail, 0, len(cFailures)*2)

	ruleCreatorArgs := ruleCreatorEmailGrouping{}

	// For each pair of user and the rule the user failed.
	for user, rules := range cFailures {
		var (
			textBuilder strings.Builder
			htmlBuilder strings.Builder
			args        = userEmailArgs{user, rules}
		)

		for _, rule := range rules {
			if _, ok := ruleCreatorArgs[rule.CreatorID]; !ok {
				ruleCreatorArgs[rule.CreatorID] = ruleCreatorEmailArgs{
					rule.CreatorName,
					rule.CreatorEmail,
					map[ruleKey][]userKey{},
				}
			}

			key := ruleKey{
				rule.ID,
				rule.Title,
				rule.Description,
				rule.ClassCode,
				rule.ClassYear,
				rule.ClassSemester,
			}
			if _, ok := ruleCreatorArgs[rule.CreatorID].RuleAndUsers[key]; !ok {
				ruleCreatorArgs[rule.CreatorID].RuleAndUsers[key] = []userKey{}
			}

			ruleCreatorArgs[rule.CreatorID].RuleAndUsers[key] = append(
				ruleCreatorArgs[rule.CreatorID].RuleAndUsers[key],
				user,
			)
		}

		if err := userTextEmail.Execute(&textBuilder, args); err != nil {
			return nil, err
		}

		if err := userHtmlEmail.Execute(&htmlBuilder, args); err != nil {
			return nil, err
		}

		mail := azmail.NewMail()
		mail.Recipients = azmail.MailRecipients{
			To: []azmail.MailAddress{{user.Email, user.Name}},
		}
		mail.Content = azmail.MailContent{
			Subject:   userEmailSubject,
			PlainText: textBuilder.String(),
			Html:      htmlBuilder.String(),
		}

		mails = append(mails, mail)
	}

	for _, args := range ruleCreatorArgs {
		var (
			textBuilder strings.Builder
			htmlBuilder strings.Builder
		)

		if err := ruleCreatorTextEmail.Execute(&textBuilder, args); err != nil {
			return nil, err
		}

		if err := ruleCreatorHtmlEmail.Execute(&htmlBuilder, args); err != nil {
			return nil, err
		}

		mail := azmail.NewMail()
		mail.Recipients = azmail.MailRecipients{
			To: []azmail.MailAddress{{args.CreatorEmail, args.CreatorName}},
		}
		mail.Content = azmail.MailContent{
			Subject:   ruleCreatorEmailSubject,
			PlainText: textBuilder.String(),
			Html:      htmlBuilder.String(),
		}

		mails = append(mails, mail)
	}

	return mails, nil
}
