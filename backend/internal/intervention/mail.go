package intervention

import (
	"strings"

	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	"github.com/darylhjd/oams/backend/pkg/azmail"
)

func (s *Service) generateMails(failures map[userKey][]model.ClassAttendanceRule) ([]*azmail.Mail, error) {
	mails := make([]*azmail.Mail, 0, len(failures))

	for user, rules := range failures {
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
