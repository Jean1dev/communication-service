package application

import (
	"communication-service/infra/config"
	"context"
	"encoding/json"
	"log"
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

func InsertNewNotification(description string, user string) error {
	db := config.GetDB()
	notification := NewNotification(description, user)
	if err := db.Insert(notification, "notifications"); err != nil {
		return err
	}

	return nil
}

func GetMyNotifications(user string) []Notification {
	db := config.GetDB()

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
