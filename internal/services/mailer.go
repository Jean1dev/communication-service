package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"mime/quotedprintable"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"time"

	"github.com/Jean1dev/communication-service/internal/dto"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/mailgun/mailgun-go/v4"
)

func downloadAttachment(attachment string) ([]byte, error) {
	resp, err := http.Get(attachment)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer requisição HTTP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erro ao baixar arquivo: status %d", resp.StatusCode)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler conteúdo do arquivo: %w", err)
	}

	return buf.Bytes(), nil
}

func sendWithMailgun(subject, recipient, htmlTemplate, attachment string) {
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

	if attachment != "" {
		bufferAttchament, err := downloadAttachment(attachment)
		if err == nil {
			message.AddBufferAttachment("anexo.pdf", bufferAttchament)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		log.Panicf("Nao foi possivel enviar o email %s", err.Error())
	}

	log.Printf("ID: %s Resp: %s\n", id, resp)
}

func sendEmailWithAttachmentSES(subject, recipient, htmlTemplate, attachment, source string) {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})

	svc := ses.New(sess)

	var emailRaw bytes.Buffer
	writer := multipart.NewWriter(&emailRaw)

	var fromName string
	if source == "notificacao@meconectei.com.br" {
		fromName = "Me Conectei"
	} else {
		fromName = "JeanLuca"
	}

	emailRaw.WriteString(fmt.Sprintf("From: %s <%s>\n", fromName, source))
	emailRaw.WriteString(fmt.Sprintf("To: %s\n", recipient))
	emailRaw.WriteString(fmt.Sprintf("Subject: %s\n", subject))
	emailRaw.WriteString("MIME-Version: 1.0\n")
	emailRaw.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\n\n", writer.Boundary()))

	mimeHeaders := textproto.MIMEHeader{}
	mimeHeaders.Set("Content-Type", "text/html; charset=UTF-8")
	mimeHeaders.Set("Content-Transfer-Encoding", "quoted-printable")

	part, _ := writer.CreatePart(mimeHeaders)
	qpWriter := quotedprintable.NewWriter(part)
	qpWriter.Write([]byte(htmlTemplate))
	qpWriter.Close()

	bufferAttchment, err := downloadAttachment(attachment)
	if err != nil {
		log.Printf("Erro ao baixar o anexo: %s", err.Error())
		return
	}

	encodedFile := base64.StdEncoding.EncodeToString(bufferAttchment)

	mimeHeaders = textproto.MIMEHeader{}
	mimeHeaders.Set("Content-Type", fmt.Sprintf("application/octet-stream; name=%s", filepath.Base("anexo.pdf")))
	mimeHeaders.Set("Content-Transfer-Encoding", "base64")
	mimeHeaders.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base("anexo.pdf")))

	part, _ = writer.CreatePart(mimeHeaders)
	part.Write([]byte(encodedFile))

	writer.Close()

	rawMessage := &ses.RawMessage{
		Data: emailRaw.Bytes(),
	}

	input := &ses.SendRawEmailInput{
		RawMessage: rawMessage,
	}
	_, err = svc.SendRawEmail(input)
	if err != nil {
		log.Printf("Não foi possível enviar o e-mail: %s", err.Error())
	} else {
		log.Println("E-mail enviado com sucesso!")
	}
}

func sendWithSES(subject, recipient, htmlTemplate, attachment, source string) {
	if attachment != "" {
		sendEmailWithAttachmentSES(subject, recipient, htmlTemplate, attachment, source)
		return
	}

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
		Source: aws.String(source),
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

	var source string
	if input.TemplateCode == 3 {
		source = "notificacao@meconectei.com.br"
	} else {
		source = "jeanlucafp@gmail.com"
	}

	if input.Recipient == "jeanlucafp@gmail.com" {
		go sendWithMailgun(input.Subject, input.Recipient, input.GetTemplate(), input.AttachmentLink)
	} else {
		go sendWithSES(input.Subject, input.Recipient, input.GetTemplate(), input.AttachmentLink, source)
	}

	return nil
}
