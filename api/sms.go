package api

import (
	"encoding/json"
	"net/http"

	"github.com/Jean1dev/communication-service/internal/services"
)

func SMSBalanceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	balance, err := services.GetAccountBalance()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(balance)
}
