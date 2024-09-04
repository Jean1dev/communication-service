package services

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Jean1dev/communication-service/internal/dto"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/mailgun/mailgun-go/v4"
)

func sendWithMailgun(subject, recipient, htmlTemplate string) {
	privateAPIKey := os.Getenv("MAILGUN_KEY")
	if privateAPIKey == "" {
		log.Print("MAILGUN_KEY not configured")
		return
	}

	domain := "central.binnoapp.com"
	mg := mailgun.NewMailgun(domain, privateAPIKey)
	sender := "Binno apps <equipe@central.binnoapp.com>"

	message := mg.NewMessage(sender, subject, "", recipient)
	message.SetHtml(htmlTemplate)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		log.Panicf("Nao foi possivel enviar o email %s", err.Error())
	}

	log.Printf("ID: %s Resp: %s\n", id, resp)
}

func sendWithSES(subject, recipient, htmlTemplate string) {
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
					Data: aws.String(htmlTemplate),
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

func AsyncSend(input dto.MailSenderInputDto) error {
	if err := input.Validate(); err != nil {
		return err
	}

	if input.Recipient == "jeanlucafp@gmail.com" {
		go sendWithMailgun(input.Subject, input.Recipient, input.GetTemplate())
	} else {
		go sendWithSES(input.Subject, input.Recipient, input.GetTemplate())
	}

	return nil
}
