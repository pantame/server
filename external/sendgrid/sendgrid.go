package sendgrid

import (
	"github.com/pantame/server/config"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendEmail(fromName, fromAddress, subject, toName, toAddress, plainTextContent, htmlContent string) error {
	from := mail.NewEmail(fromName, fromAddress)
	to := mail.NewEmail(toName, toAddress)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(config.AccessTokens().SendGridKey)

	_, err := client.Send(message)
	if err != nil {
		return err
	}
	return nil
}
