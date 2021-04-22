package mail

import (
	"github.com/pantame/server/external/sendgrid"
	"net/mail"
)

type Mail interface {
	SendMessage() error
}

type Header struct {
	From    mail.Address
	Subject string
	To      mail.Address
}

type Message struct {
	Header    Header
	Content   string
	PlainText bool
}

func NewEmail(header Header, content string, plainText bool) Mail {
	return &Message{
		Header: Header{
			From: mail.Address{
				Name:    header.From.Name,
				Address: header.From.Address,
			},
			Subject: header.Subject,
			To: mail.Address{
				Name:    header.To.Name,
				Address: header.To.Address,
			},
		},
		PlainText: plainText,
		Content:   content,
	}
}

func (msg Message) SendMessage() error {
	var plainTextContent, htmlContent string

	if msg.PlainText {
		plainTextContent = msg.Content
	} else {
		htmlContent = msg.Content
	}

	return sendgrid.SendEmail(msg.Header.From.Name, msg.Header.From.Address, msg.Header.Subject, msg.Header.To.Name, msg.Header.To.Address, plainTextContent, htmlContent)
}
