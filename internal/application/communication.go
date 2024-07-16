package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Jean1dev/communication-service/configs"
	"github.com/Jean1dev/communication-service/internal/services"
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

	fmt.Println("Resposta:", string(body))

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)

	if err != nil {
		log.Println("Erro ao fazer o unmarshal do JSON:", err)
		return "", err
	}

	phoneNumber, ok := data["phoneNumber"].([]interface{})
	if !ok {
		log.Println("Erro ao acessar o campo 'phoneNumber'")
		return "", errors.New("Erro ao acessar o campo 'phoneNumber'")
	}

	return phoneNumber[0].(string), nil
}

func sendEmail(recipient string, message string) {
	services.AsyncSend("Notificacao", message, recipient)
}

func sendSMS(message string, recipient string) {
	userPhone, err := searchForPhone(recipient)
	if err != nil {
		return
	}

	services.DispatchSMS(userPhone, message)
}

func SendCommunications(message string, recipients []string, types []string) {
	for _, typeCommunication := range types {
		if typeCommunication == "email" {
			for _, recipient := range recipients {
				go sendEmail(recipient, message)
			}
		}

		if typeCommunication == "sms" {
			for _, recipient := range recipients {
				go sendSMS(message, recipient)
			}
		}
	}
}
