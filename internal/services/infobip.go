package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Jean1dev/communication-service/internal/infra/database"
)

type MessagePayload struct {
	Messages []Message `json:"messages"`
}

type Message struct {
	Destinations []Destination `json:"destinations"`
	From         string        `json:"from"`
	Text         string        `json:"text"`
}

type Destination struct {
	To string `json:"to"`
}

func numberFormat(number string) string {
	if !strings.HasPrefix(number, "55") {
		number = "55" + number
	}

	log.Printf("Number formatted: %s", number)
	return number
}

func buildPayload(to []string, text string) MessagePayload {
	destinations := make([]Destination, len(to))
	for i, number := range to {
		destinations[i] = Destination{To: numberFormat(number)}
	}

	message := Message{
		Destinations: destinations,
		From:         "ServiceSMS",
		Text:         text,
	}

	return MessagePayload{Messages: []Message{message}}
}

func getMyNumber() []string {
	return []string{os.Getenv("MY_NUMBER")}
}

func DispatchSMS(to []string, text string) {
	url := "https://vj44dv.api.infobip.com/sms/2/text/advanced"
	method := "POST"

	if len(to) == 0 {
		to = getMyNumber()
	}

	messagePayload := buildPayload(to, text)
	payload, err := json.Marshal(messagePayload)
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(string(payload)))

	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Add("Authorization", os.Getenv("INFOBIP_KEY"))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	go auditMessage(string(body))
}

func auditMessage(messageStr string) {
	db := database.GetDB()

	document := map[string]interface{}{
		"result": messageStr,
	}

	if err := db.Insert(document, "infobip_audit"); err != nil {
		log.Printf("Error inserting message into database: %v", err)
		return
	}

	log.Printf("Message inserted into database: %s", messageStr)
}
