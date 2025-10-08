package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Jean1dev/communication-service/internal/application"
	"github.com/Jean1dev/communication-service/internal/dto"
	"github.com/Jean1dev/communication-service/internal/infra/database"
)

var alertService *application.AlertService

func initAlertService() {
	if alertService == nil {
		db := database.GetDB()
		alertService = application.NewAlertService(db)
	}
}

func AlertHandler(w http.ResponseWriter, r *http.Request) {
	initAlertService()

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		handleCreateAlert(w, r)
	case http.MethodGet:
		handleGetAlerts(w, r)
	case http.MethodPut:
		handleUpdateAlert(w, r)
	case http.MethodDelete:
		handleDeleteAlert(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleCreateAlert(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateAlertInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	alert, err := alertService.CreateAlert(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := alert.ToResponse()
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func handleGetAlerts(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	alertID := strings.TrimPrefix(path, "/alerts/")

	if alertID != "" && alertID != "/alerts" {
		alert, err := alertService.GetAlertByID(alertID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		response := alert.ToResponse()
		json.NewEncoder(w).Encode(response)
		return
	}

	userEmail := r.URL.Query().Get("user_email")
	if userEmail == "" {
		http.Error(w, "user_email query parameter is required", http.StatusBadRequest)
		return
	}

	alerts, err := alertService.GetAlertsByUserEmail(userEmail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var responses []dto.AlertResponse
	for _, alert := range alerts {
		responses = append(responses, alert.ToResponse())
	}

	json.NewEncoder(w).Encode(responses)
}

func handleUpdateAlert(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	alertID := strings.TrimPrefix(path, "/alerts/")

	if alertID == "" || alertID == "/alerts" {
		http.Error(w, "alert id is required in path", http.StatusBadRequest)
		return
	}

	var input dto.UpdateAlertInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	alert, err := alertService.UpdateAlert(alertID, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := alert.ToResponse()
	json.NewEncoder(w).Encode(response)
}

func handleDeleteAlert(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	alertID := strings.TrimPrefix(path, "/alerts/")

	if alertID == "" || alertID == "/alerts" {
		http.Error(w, "alert id is required in path", http.StatusBadRequest)
		return
	}

	err := alertService.DeleteAlert(alertID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func AlertToggleStatusHandler(w http.ResponseWriter, r *http.Request) {
	initAlertService()

	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Path
	alertID := strings.TrimPrefix(path, "/alerts/toggle/")

	if alertID == "" || alertID == path {
		http.Error(w, "alert id is required in path", http.StatusBadRequest)
		return
	}

	alert, err := alertService.ToggleAlertStatus(alertID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := alert.ToResponse()
	json.NewEncoder(w).Encode(response)
}
