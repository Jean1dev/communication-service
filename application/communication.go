package application

import (
	"communication-service/services"
	"log"
)

func sendEmail(recipient string, message string) {
	services.AsyncSend("Notificacao", message, recipient)
}

func sendSMS(message string, recipient string) {
	log.Println("SMS enviado para " + recipient)
}

func SendCommunications(message string, recipients []string, types []string) {
	for _, typeCommunication := range types {
		if typeCommunication == "email" {
			for _, recipient := range recipients {
				sendEmail(recipient, message)
			}
		}

		if typeCommunication == "sms" {
			for _, recipient := range recipients {
				sendSMS(message, recipient)
			}
		}
	}
}
