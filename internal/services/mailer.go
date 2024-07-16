package services

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/mailgun/mailgun-go/v4"
)

func sendWithMailgun(subject string, body string, recipient string) {
	privateAPIKey := os.Getenv("MAILGUN_KEY")
	if privateAPIKey == "" {
		log.Print("MAILGUN_KEY not configured")
		return
	}

	domain := "central.binnoapp.com"
	mg := mailgun.NewMailgun(domain, privateAPIKey)
	sender := "Binno apps <equipe@central.binnoapp.com>"

	message := mg.NewMessage(sender, subject, body, recipient)
	message.SetHtml(Default(subject, body))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		log.Panicf("Nao foi possivel enviar o email %s", err.Error())
	}

	log.Printf("ID: %s Resp: %s\n", id, resp)
}

func sendWithSES(subject string, body string, recipient string) {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})

	svc := ses.New(sess)

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Data: aws.String(Default(subject, body)),
				},
			},
			Subject: &ses.Content{
				Data: aws.String(subject),
			},
		},
		Source: aws.String("jeanlucafp@gmail.com"),
	}

	_, err := svc.SendEmail(input)
	if err != nil {
		log.Printf("Nao foi possivel enviar o email %s", err.Error())
	}
}

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

	if recipient == "jeanlucafp@gmail.com" {
		go sendWithMailgun(subject, body, recipient)
	} else {
		go sendWithSES(subject, body, recipient)
	}

	return nil
}
