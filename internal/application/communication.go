package application

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Jean1dev/communication-service/configs"
	"github.com/Jean1dev/communication-service/internal/dto"
	"github.com/Jean1dev/communication-service/internal/services"
	"github.com/valyala/fastjson"
)

func searchForPhone(recipient string) (string, error) {
	url := fmt.Sprintf("%s/get-user-data?email=%s", configs.CaixinhaServer, recipient)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Erro ao criar a requisição para a Caixinha %s\n", err)
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()

	req = req.WithContext(ctx)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("Erro ao fazer a requisição para a Caixinha %s\n", err)
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Print("Erro ao ler o corpo da resposta:", err)
		return "", err
	}

	var p fastjson.Parser
	v, err := p.Parse(string(body))
	if err != nil {
		log.Println("Erro ao acessar o campo 'phoneNumber'")
		return "", errors.New("Erro ao acessar o campo 'phoneNumber'")
	}

	phoneNumber := string(v.GetStringBytes("phoneNumber"))
	return phoneNumber, nil
}

func sendEmail(input dto.MailSenderInputDto) {
	services.AsyncSend(input)
}

func sendSMS(message string, recipient []string) {
	userPhones := make([]string, 0)
	for _, recipient := range recipient {
		if phone, err := strconv.Atoi(recipient); err == nil {
			userPhones = append(userPhones, strconv.Itoa(phone))
		} else if phone, err := searchForPhone(recipient); err == nil {
			userPhones = append(userPhones, phone)
		}
	}

	services.DispatchSMS(userPhones, message)
}

func SendCommunications(input dto.MailSenderInputDto, recipients, types []string) {
	for _, typeCommunication := range types {
		if typeCommunication == "email" {
			for _, recipient := range recipients {
				input.Recipient = recipient
				go sendEmail(input)
			}
		}

		if typeCommunication == "sms" {
			go sendSMS(input.Body, recipients)
		}
	}
}
