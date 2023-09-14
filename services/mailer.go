package services

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

var domain = "central.binnoapp.com"

func AsyncSend(subject string, body string, recipient string) error {
	if subject == "" {
		return errors.New("Subject cannot be empty")
	}
	if body == "" {
		return errors.New("body cannot be empty")
	}
	if recipient == "" {
		return errors.New("recipient cannot be empty")
	}

	privateAPIKey := os.Getenv("MAILGUN_KEY")
	if privateAPIKey == "" {
		log.Print("MAILGUN_KEY not configured")
		return nil
	}

	go func() {
		mg := mailgun.NewMailgun(domain, privateAPIKey)
		sender := "Binno apps <equipe@central.binnoapp.com>"

		message := mg.NewMessage(sender, subject, body, recipient)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		resp, id, err := mg.Send(ctx, message)

		if err != nil {
			log.Panicf("Nao foi possivel enviar o email %s", err.Error())
		}

		log.Printf("ID: %s Resp: %s\n", id, resp)
	}()

	return nil
}
