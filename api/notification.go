package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Jean1dev/communication-service/configs"
	"github.com/Jean1dev/communication-service/internal/application"
	"github.com/Jean1dev/communication-service/internal/dto"
)

type NotificationPost struct {
	Desc         string   `json:"desc"`
	User         string   `json:"user"`
	Caixinhas    []string `json:"caixinhas"`
	Comunicacoes []string `json:"types"`
}

type NotificationSMS struct {
	Desc       string   `json:"desc"`
	Recipients []string `json:"recipients"`
}

type MarkNotificationAsRead struct {
	User string   `json:"user"`
	All  bool     `json:"all"`
	Ids  []string `json:"ids"`
}

func NotificationHandler(w http.ResponseWriter, r *http.Request) {
	configs.AllowAllOrigins(w, r)
	method := r.Method
	if method == "POST" {

		if strings.HasPrefix(r.URL.Path, "/notificacao/mark-as-read") {
			markAllAsRead(w, r)
			return
		}

		if strings.HasPrefix(r.URL.Path, "/notificacao/sms") {
			sendSMS(w, r)
			return
		}

		doPost(w, r)
		return
	}

	http.Error(w, "Method not allowed", http.StatusBadRequest)
}

func markAllAsRead(w http.ResponseWriter, r *http.Request) {
	var payload MarkNotificationAsRead
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if payload.All {
		if err := application.MarkNotificationAsRead(make([]string, 0), payload.User); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if err := application.MarkNotificationAsRead(payload.Ids, payload.User); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func doPost(w http.ResponseWriter, r *http.Request) {
	var payload NotificationPost
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var err error
	users := []string{payload.User}
	if len(payload.Caixinhas) > 0 {
		users = application.NewNofiticationForCaixinha(payload.Desc, payload.User, payload.Caixinhas)
	} else {
		err = application.InsertNewNotification(payload.Desc, payload.User)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	communicationInput := dto.MailSenderInputDto{
		Body: payload.Desc,
	}
	application.SendCommunications(communicationInput, users, payload.Comunicacoes)
}

func removeSpecialCaracteres(numbers []string) []string {
	for i, number := range numbers {
		number = strings.ReplaceAll(number, "(", "")
		number = strings.ReplaceAll(number, ")", "")
		number = strings.ReplaceAll(number, "-", "")
		number = strings.ReplaceAll(number, " ", "")
		numbers[i] = number
	}
	return numbers
}

func sendSMS(w http.ResponseWriter, r *http.Request) {
	var payload NotificationSMS
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	communicationInput := dto.MailSenderInputDto{
		Body: payload.Desc,
	}

	application.SendCommunications(communicationInput, removeSpecialCaracteres(payload.Recipients), []string{"sms"})
	w.WriteHeader(http.StatusOK)
}
