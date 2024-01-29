package services_test

import (
	"communication-service/services"
	"os"
	"testing"
)

func TestSendMail(t *testing.T) {
	subject := "Fancy subject!"
	body := "Hello from Mailgun Go!"
	recipient := "jeanlucafp@gmail.com"

	os.Setenv("MAILGUN_KEY", "chave-mock")

	services.AsyncSend(subject, body, recipient)
}
