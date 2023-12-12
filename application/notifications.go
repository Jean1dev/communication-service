package application

import (
	"communication-service/infra/database"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Notification struct {
	Id          string `json:"id"`
	CreatedAt   string `json:"createdAt"`
	Description string `json:"description"`
	Read        bool   `json:"read"`
	User        string `json:"user"`
}

func NewNotification(description string, user string) *Notification {
	return &Notification{
		CreatedAt:   time.Now().String(),
		Description: description,
		Read:        false,
		User:        user,
	}
}

func verifyIfElementIsOnTheList(list []string, element string) bool {
	for _, it := range list {
		if it == element {
			return true
		}
	}

	return false
}

func NewNofiticationForCaixinha(description string, user string, caixinhasId []string) []string {
	notificacoesList := []string{user}
	for _, it := range caixinhasId {
		url := fmt.Sprintf("https://emprestimo-caixinha.azurewebsites.net/api/dados-analise?caixinhaId=%s", it)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Printf("Erro ao criar a requisição para a Caixinha %s: %s\n", it, err)
			continue
		}

		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			fmt.Printf("Erro ao fazer a requisição para a Caixinha %s: %s\n", it, err)
			continue
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Print("Erro ao ler o corpo da resposta:", err)
			continue
		}

		fmt.Println("Resposta:", string(body))

		var data map[string]interface{}
		err = json.Unmarshal(body, &data)

		if err != nil {
			log.Println("Erro ao fazer o unmarshal do JSON:", err)
			continue
		}

		membros, ok := data["membros"].([]interface{})
		if !ok {
			log.Println("Erro ao acessar o campo 'membros'")
			continue
		}

		for _, membro := range membros {
			membroMap, ok := membro.(map[string]interface{})
			if !ok {
				log.Println("Erro ao converter movimentacao para map[string]interface{}")
				continue
			}

			name := membroMap["name"].(string)
			email := membroMap["email"].(string)

			log.Printf("elemento %s email %s", name, email)

			if verifyIfElementIsOnTheList(notificacoesList, email) {
				continue
			}

			notificacoesList = append(notificacoesList, email)
		}
	}

	for _, it := range notificacoesList {
		if err := InsertNewNotification(description, it); err != nil {
			log.Fatalf(err.Error())
		}
	}

	return notificacoesList
}

func InsertNewNotification(description string, user string) error {
	db := database.GetDB()
	notification := NewNotification(description, user)
	if err := db.Insert(notification, "notifications"); err != nil {
		return err
	}

	return nil
}

func GetMyNotifications(user string) []Notification {
	db := database.GetDB()

	err, cursor := db.FindAll("notifications", bson.D{{Key: "user", Value: user}}, options.Find())
	if err != nil {
		panic(err)
	}

	var results []Notification
	for cursor.Next(context.Background()) {
		var doc bson.M
		err := cursor.Decode(&doc)
		if err != nil {
			log.Fatal(err)
		}

		jsonData, err := bson.MarshalExtJSON(doc, false, false)
		if err != nil {
			log.Fatal(err)
		}

		var notification Notification
		err = json.Unmarshal(jsonData, &notification)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, notification)
	}
	return results
}
