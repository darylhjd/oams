package main

import (
	"log"
	"os"

	"github.com/darylhjd/oams/backend/pkg/azmail"
)

func main() {
	c, err := azmail.NewClient(
		os.Getenv("AZURE_EMAIL_ENDPOINT"),
		os.Getenv("AZURE_EMAIL_ACCESS_KEY"),
		os.Getenv("AZURE_EMAIL_SENDER_ADDRESS"),
	)
	if err != nil {
		log.Fatal(err)
	}

	mail := azmail.NewMail()
	mail.Recipients = azmail.MailRecipients{
		To: []azmail.MailAddress{
			{
				Address:     "",
				DisplayName: "",
			},
		},
	}
	mail.Content = azmail.MailContent{
		Subject:   "Azure Test Email",
		PlainText: "This is the test email content",
	}

	if err = c.SendMails(mail); err != nil {
		log.Fatal(err)
	}

	log.Println("email sending successful")
}
