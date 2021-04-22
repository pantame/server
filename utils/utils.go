package utils

import (
	"fmt"
	"github.com/pantame/server/config"
	"github.com/pantame/server/utils/mail"
	"math/rand"
	netMail "net/mail"
	"time"
)

func New6DigitCode() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%d%d%d%d%d%d", r.Intn(9) + 1, r.Intn(9) + 1, r.Intn(9) + 1, r.Intn(9) + 1, r.Intn(9) + 1, r.Intn(9) + 1)
}

func SendEmailVerificationCode(name, address string) (string, error) {
	code := New6DigitCode()
	header := mail.Header{
		From:    config.Emails().FromNoReply,
		Subject: "Código de Verificação - Pantame",
		To: netMail.Address{
			Name:    name,
			Address: address,
		},
	}
	message := mail.NewEmail(header, fmt.Sprintf("Seu código de verificação Pantame é: %s", code), true)
	return code, message.SendMessage()
}
