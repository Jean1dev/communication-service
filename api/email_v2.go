package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Jean1dev/communication-service/internal/dto"
	"github.com/Jean1dev/communication-service/internal/infra/database"
	"github.com/Jean1dev/communication-service/internal/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SendEmailV2Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input dto.SendEmailV2Dto
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := services.SendEmailV2(input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db := database.GetDB()
	emailEntity := struct {
		To        string             `json:"to"`
		CreatedAt primitive.DateTime `bson:"createdAt" json:"createdAt"`
	}{
		To:        input.To,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}
	go db.Insert(emailEntity, email_collection)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "pending",
		"message": "email queued successfully",
	})
}

func EmailTemplatesV2Handler(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/v2/email/templates")
	slug = strings.TrimPrefix(slug, "/")

	switch {
	case slug == "" && r.Method == http.MethodGet:
		listTemplates(w, r)
	case slug == "" && r.Method == http.MethodPost:
		createTemplate(w, r)
	case slug != "" && r.Method == http.MethodGet:
		getTemplate(w, r, slug)
	case slug != "" && r.Method == http.MethodPut:
		updateTemplate(w, r, slug)
	case slug != "" && r.Method == http.MethodDelete:
		deleteTemplate(w, r, slug)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func createTemplate(w http.ResponseWriter, r *http.Request) {
	var t dto.EmailTemplateDto
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := services.CreateEmailTemplate(t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "created",
		"message": "template created successfully",
		"slug":    t.Slug,
	})
}

func listTemplates(w http.ResponseWriter, r *http.Request) {
	templates, err := services.ListEmailTemplates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(templates)
}

func getTemplate(w http.ResponseWriter, r *http.Request, slug string) {
	t, err := services.GetEmailTemplate(slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func updateTemplate(w http.ResponseWriter, r *http.Request, slug string) {
	var t dto.EmailTemplateDto
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := services.UpdateEmailTemplate(slug, t); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "updated",
		"message": "template updated successfully",
	})
}

func deleteTemplate(w http.ResponseWriter, r *http.Request, slug string) {
	if err := services.DeleteEmailTemplate(slug); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
