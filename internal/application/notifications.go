package application

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Jean1dev/communication-service/configs"
	"github.com/Jean1dev/communication-service/internal/infra/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Notification struct {
	Id          string `json:"id"`
	CreatedAt   string `json:"createdAt"`
	Description string `json:"description"`
	Read        bool   `json:"read"`
	User        string `json:"user"`
	Type        string `json:"type"`
}

func NewNotification(description string, user string) *Notification {
	return &Notification{
		CreatedAt:   time.Now().String(),
		Description: description,
		Read:        false,
		User:        user,
		Type:        "user_info",
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

func MarkNotificationAsRead(ids []string, user string) error {
	db := database.GetDB()

	for id := range ids {
		filter := bson.D{{Key: "_id", Value: id}}
		update := bson.D{
			{Key: "$set", Value: bson.D{{Key: "read", Value: true}}},
		}

		if err := db.UpdateOne("notifications", filter, update); err != nil {
			return err
		}
	}

	if len(ids) == 0 {
		filter := bson.D{{Key: "user", Value: user}}
		update := bson.D{
			{Key: "$set", Value: bson.D{{Key: "read", Value: true}}},
		}

		if err := db.UpdateMany("notifications", filter, update); err != nil {
			return err
		}
	}

	return nil
}

func NewNofiticationForCaixinha(description string, user string, caixinhasId []string) []string {
	notificacoesList := []string{user}
	for _, it := range caixinhasId {
		url := fmt.Sprintf("%s/dados-analise?caixinhaId=%s", configs.CaixinhaServer, it)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Printf("Erro ao criar a requisição para a Caixinha %s: %s\n", it, err)
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
		defer cancel()

		req = req.WithContext(ctx)

		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			fmt.Printf("Erro ao fazer a requisição para a Caixinha %s: %s\n", it, err)
			continue
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
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
		if err := newAlertNotification(description, it); err != nil {
			log.Fatalf(err.Error())
		}
	}

	return notificacoesList
}

func newAlertNotification(description string, user string) error {
	db := database.GetDB()
	notification := NewNotification(description, user)
	notification.Type = "new_feature"
	if err := db.Insert(notification, "notifications"); err != nil {
		return err
	}

	return nil
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

		if notification.Read {
			continue
		}

		notification.Id = doc["_id"].(primitive.ObjectID).Hex()
		results = append(results, notification)
	}
	return results
}
