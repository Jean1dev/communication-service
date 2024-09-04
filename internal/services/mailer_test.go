package services_test

import (
	"os"
	"testing"

	"github.com/Jean1dev/communication-service/internal/dto"
	"github.com/Jean1dev/communication-service/internal/services"
)

func TestSendMail(t *testing.T) {
	subject := "Fancy subject!"
	body := "Hello from Mailgun Go!"
	recipient := "jeanlucafp@gmail.com"
	input := dto.MailSenderInputDto{
		Subject:   subject,
		Body:      body,
		Recipient: recipient,
	}

	os.Setenv("MAILGUN_KEY", "chave-mock")

	services.AsyncSend(input)
}
