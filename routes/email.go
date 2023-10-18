package routes

import (
	"communication-service/infra/database"
	"communication-service/services"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

var email_collection = "emails_sending"

type Email struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func EmailEstatisticasHandler(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB()
	umaSemanaAtras := time.Now().AddDate(0, 0, -7)
	filter := bson.D{{Key: "createdat", Value: bson.D{{Key: "$lt", Value: umaSemanaAtras}}}}

	qtd, err := db.CountDocuments(email_collection, filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = services.AsyncSend("Quantidade de emails enviado na semanda", fmt.Sprintf("a quantidade foi %d", qtd), "jeanlucafp@gmail.com")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func EmailHandler(w http.ResponseWriter, r *http.Request) {
	var emailData Email
	err := json.NewDecoder(r.Body).Decode(&emailData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = services.AsyncSend(emailData.Subject, emailData.Message, emailData.To)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	db := database.GetDB()

	emailEntity := struct {
		To        string    `json:"to"`
		CreatedAt time.Time `json:"createdAt"`
	}{
		CreatedAt: time.Now(),
		To:        emailData.To,
	}

	go db.Insert(emailEntity, email_collection)

	response := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}{
		Status:  "pending",
		Message: "sending email solicited sucessfuly",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
