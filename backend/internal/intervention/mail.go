package intervention

import (
	"strings"

	"github.com/darylhjd/oams/backend/pkg/azmail"
)

func (s *Service) generateMails(cFailures checkFailures) ([]*azmail.Mail, error) {
	mails := make([]*azmail.Mail, 0, len(cFailures))

	for user, rules := range cFailures {
		var (
			textBuilder strings.Builder
			htmlBuilder strings.Builder
			args        = EmailArgs{user, rules}
		)

		if err := textEmail.Execute(&textBuilder, args); err != nil {
			return nil, err
		}

		if err := htmlEmail.Execute(&htmlBuilder, args); err != nil {
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

	return mails, nil
}
