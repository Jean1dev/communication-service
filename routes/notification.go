package routes

import (
	"communication-service/application"
	"communication-service/infra/config"
	"encoding/json"
	"net/http"
	"strings"
)

type NotificationPost struct {
	Desc      string   `json:"desc"`
	User      string   `json:"user"`
	Caixinhas []string `json:"caixinhas"`
}

type MarkNotificationAsRead struct {
	User string   `json:"user"`
	All  bool     `json:"all"`
	Ids  []string `json:"ids"`
}

func NotificationHandler(w http.ResponseWriter, r *http.Request) {
	config.AllowAllOrigins(w, r)
	method := r.Method
	if method == "POST" {

		if strings.HasPrefix(r.URL.Path, "/notificacao/mark-as-read") {
			markAllAsRead(w, r)
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
	if len(payload.Caixinhas) > 1 {
		application.NewNofiticationForCaixinha(payload.Desc, payload.User, payload.Caixinhas)
	} else {
		err = application.InsertNewNotification(payload.Desc, payload.User)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
