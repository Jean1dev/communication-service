package routes

import (
	"communication-service/application"
	"communication-service/infra/config"
	"encoding/json"
	"net/http"
)

type NotificationPost struct {
	Desc      string   `json:"desc"`
	User      string   `json:"user"`
	Caixinhas []string `json:"caixinhas"`
}

func NotificationHandler(w http.ResponseWriter, r *http.Request) {
	config.AllowAllOrigins(w, r)
	method := r.Method
	if method == "POST" {
		doPost(w, r)
		return
	}

	http.Error(w, "Method not allowed", http.StatusBadRequest)
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
