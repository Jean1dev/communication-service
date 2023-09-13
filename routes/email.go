package routes

import (
	"encoding/json"
	"net/http"
)

type Email struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func EmailHandler(w http.ResponseWriter, r *http.Request) {
	var emailData Email
	err := json.NewDecoder(r.Body).Decode(&emailData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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
