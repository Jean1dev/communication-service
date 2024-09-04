package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Jean1dev/communication-service/internal/dto"
	"github.com/Jean1dev/communication-service/internal/infra/database"
	"github.com/Jean1dev/communication-service/internal/services"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var email_collection = "emails_sending"

func EmailEstatisticasHandler(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB()
	umaSemanaAtras := time.Now().AddDate(0, 0, -7)
	umaSemanaAtrasMongo := primitive.NewDateTimeFromTime(umaSemanaAtras)
	filter := bson.D{{Key: "createdAt", Value: bson.D{{Key: "$gt", Value: umaSemanaAtrasMongo}}}}

	qtd, err := db.CountDocuments(email_collection, filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mailInput := dto.MailSenderInputDto{
		Subject:   "Quantidade de emails enviado na semanda",
		Body:      fmt.Sprintf("a quantidade foi %d", qtd),
		Recipient: "jeanlucafp@gmail.com",
	}
	err = services.AsyncSend(mailInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func EmailHandler(w http.ResponseWriter, r *http.Request) {
	var emailData dto.MailSenderInputDto
	err := json.NewDecoder(r.Body).Decode(&emailData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = services.AsyncSend(emailData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	db := database.GetDB()

	emailEntity := struct {
		To        string             `json:"to"`
		CreatedAt primitive.DateTime `bson:"createdAt" json:"createdAt"`
	}{
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		To:        emailData.Recipient,
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
