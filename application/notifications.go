package application

import (
	"communication-service/infra/config"
	"context"
	"encoding/json"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

type Notification struct {
	Id          string `json:"id"`
	CreatedAt   string `json:"createdAt"`
	Description string `json:"description"`
	Read        bool   `json:"read"`
}

func GetMyNotifications(whoNotifications string) []Notification {
	db := config.GetDB()

	err, cursor := db.FindAll("notifications", bson.D{{Key: "quem", Value: whoNotifications}})
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
