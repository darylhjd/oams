package intervention

import (
	"strings"

	"github.com/darylhjd/oams/backend/pkg/azmail"
)

func (s *Service) generateNotificationMails(users userFailedRules, ruleCreators ruleCreatorRuleFailedUsers) ([]*azmail.Mail, error) {
	mails := make([]*azmail.Mail, 0, len(users)+len(ruleCreators))

	// For each pair of user and the rule the user failed.
	for user, rules := range users {
		var (
			textBuilder strings.Builder
			htmlBuilder strings.Builder
			args        = userEmailArgs{user, rules}
		)

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

	// For each pair of rule creator and their rule with failed users.
	for creator, ruleWithFailedUsers := range ruleCreators {
		var (
			textBuilder strings.Builder
			htmlBuilder strings.Builder
			args        = ruleCreatorEmailArgs{creator, ruleWithFailedUsers}
		)

		if err := ruleCreatorTextEmail.Execute(&textBuilder, args); err != nil {
			return nil, err
		}

		if err := ruleCreatorHtmlEmail.Execute(&htmlBuilder, args); err != nil {
			return nil, err
		}

		mail := azmail.NewMail()
		mail.Recipients = azmail.MailRecipients{
			To: []azmail.MailAddress{{creator.Email, creator.Name}},
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
